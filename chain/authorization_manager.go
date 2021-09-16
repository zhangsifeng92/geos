package chain

import (
	"github.com/zhangsifeng92/geos/chain/types"
	. "github.com/zhangsifeng92/geos/chain/types/generated_containers"
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto/rlp"
	"github.com/zhangsifeng92/geos/database"
	"github.com/zhangsifeng92/geos/entity"
	. "github.com/zhangsifeng92/geos/exception"
	. "github.com/zhangsifeng92/geos/exception/try"
	"github.com/zhangsifeng92/geos/log"
)

var noopCheckTime *func()

type AuthorizationManager struct {
	control *Controller
	db      database.DataBase
}

func newAuthorizationManager(control *Controller) *AuthorizationManager {
	azInstance := &AuthorizationManager{}
	azInstance.control = control
	azInstance.db = control.DB
	return azInstance
}

type PermissionIdType common.IdType

func (a *AuthorizationManager) CreatePermission(account common.AccountName,
	name common.PermissionName,
	parent PermissionIdType,
	auth types.Authority,
	initialCreationTime common.TimePoint,
) *entity.PermissionObject {
	creationTime := initialCreationTime
	if creationTime == common.TimePoint(0) {
		creationTime = a.control.PendingBlockTime()
	}

	permUsage := entity.PermissionUsageObject{}
	permUsage.LastUsed = creationTime
	err := a.db.Insert(&permUsage)
	if err != nil {
		log.Error("CreatePermission is error: %s", err)
	}

	perm := entity.PermissionObject{
		UsageId:     permUsage.ID,
		Parent:      common.IdType(parent),
		Owner:       account,
		Name:        name,
		LastUpdated: creationTime,
		Auth:        auth.ToSharedAuthority(),
	}
	err = a.db.Insert(&perm)
	if err != nil {
		log.Error("CreatePermission is error: %s", err)
	}
	return &perm
}

func (a *AuthorizationManager) ModifyPermission(permission *entity.PermissionObject, auth *types.Authority) {
	err := a.db.Modify(permission, func(po *entity.PermissionObject) {
		po.Auth = (*auth).ToSharedAuthority()
		po.LastUpdated = a.control.PendingBlockTime()
	})
	if err != nil {
		log.Error("ModifyPermission is error: %s", err)
	}
}

func (a *AuthorizationManager) RemovePermission(permission *entity.PermissionObject) {
	index, err := a.db.GetIndex("byParent", entity.PermissionObject{})
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	itr, err := index.LowerBound(entity.PermissionObject{Parent: permission.ID})
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	EosAssert(index.CompareEnd(itr), &ActionValidateException{}, "Cannot remove a permission which has children. Remove the children first.")
	usage := entity.PermissionUsageObject{ID: permission.UsageId}
	err = a.db.Find("id", usage, &usage)
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	err = a.db.Remove(&usage)
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	err = a.db.Remove(permission)
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
}

func (a *AuthorizationManager) UpdatePermissionUsage(permission *entity.PermissionObject) {
	puo := entity.PermissionUsageObject{}
	puo.ID = permission.UsageId
	err := a.db.Find("id", puo, &puo)
	if err != nil {
		log.Error("UpdatePermissionUsage is error: %s", err)
	}
	err = a.db.Modify(&puo, func(p *entity.PermissionUsageObject) {
		puo.LastUsed = a.control.PendingBlockTime()
	})
	if err != nil {
		log.Error("UpdatePermissionUsage is error: %s", err)
	}
}

func (a *AuthorizationManager) GetPermissionLastUsed(permission *entity.PermissionObject) common.TimePoint {
	puo := entity.PermissionUsageObject{}
	puo.ID = permission.UsageId
	err := a.db.Find("id", puo, &puo)
	if err != nil {
		log.Error("GetPermissionLastUsed is error: %s", err)
	}
	return puo.LastUsed
}

func (a *AuthorizationManager) FindPermission(level *common.PermissionLevel) (p *entity.PermissionObject) {
	Try(func() {
		EosAssert(!level.Actor.Empty() && !level.Permission.Empty(), &InvalidPermission{}, "Invalid permission")
		po := entity.PermissionObject{}
		po.Owner = level.Actor
		po.Name = level.Permission
		err := a.db.Find("byOwner", po, &po)
		if err != nil {
			//log.Warn("%v@%v don't find", po.Owner, po.Name)
			p = nil
			return
		}
		p = &po
	}).EosRethrowExceptions(&PermissionQueryException{}, "Failed to retrieve permission: %v", level)
	return p
}

func (a *AuthorizationManager) GetPermission(level *common.PermissionLevel) (p *entity.PermissionObject) {
	Try(func() {
		EosAssert(!level.Actor.Empty() && !level.Permission.Empty(), &InvalidPermission{}, "Invalid permission")
		po := entity.PermissionObject{}
		po.Owner = level.Actor
		po.Name = level.Permission
		err := a.db.Find("byOwner", po, &po)
		if err != nil {
			//log.Warn("%v@%v don't find", po.Owner, po.Name)
			EosAssert(false, &PermissionQueryException{}, "Failed to retrieve permission: %v", level)
		}
		p = &po
	}).EosRethrowExceptions(&PermissionQueryException{}, "Failed to retrieve permission: %v", level)
	return p
}

func (a *AuthorizationManager) LookupLinkedPermission(authorizerAccount common.AccountName,
	scope common.AccountName,
	actName common.ActionName,
) (p *common.PermissionName) {
	Try(func() {
		link := entity.PermissionLinkObject{}
		link.Account = authorizerAccount
		link.Code = scope
		link.MessageType = actName
		err := a.db.Find("byActionName", link, &link)
		if err != nil {
			link.MessageType = common.AccountName(common.N(""))
			err = a.db.Find("byActionName", link, &link)
		}
		if err == nil {
			p = &link.RequiredPermission
			return
		}
	}).End()

	return p
}

func (a *AuthorizationManager) LookupMinimumPermission(authorizerAccount common.AccountName,
	scope common.AccountName,
	actName common.ActionName,
) (p *common.PermissionName) {
	if scope == common.DefaultConfig.SystemAccountName {
		EosAssert(actName != UpdateAuth{}.GetName() &&
			actName != DeleteAuth{}.GetName() &&
			actName != LinkAuth{}.GetName() &&
			actName != UnLinkAuth{}.GetName() &&
			actName != CancelDelay{}.GetName(),
			&UnlinkableMinPermissionAction{}, "cannot call lookup_minimum_permission on native actions that are not allowed to be linked to minimum permissions")
	}
	Try(func() {
		linkedPermission := a.LookupLinkedPermission(authorizerAccount, scope, actName)
		if common.Empty(linkedPermission) {
			p = &common.DefaultConfig.ActiveName
			return
		}

		if *linkedPermission == common.PermissionName(common.DefaultConfig.EosioAnyName) {
			return
		}

		p = linkedPermission
		return
	}).End()
	return p
}

func (a *AuthorizationManager) CheckUpdateAuthAuthorization(update UpdateAuth, auths []common.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "UpdateAuth action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == update.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")
	minPermission := a.FindPermission(&common.PermissionLevel{Actor: update.Account, Permission: update.Permission})
	if minPermission == nil {
		permission := a.GetPermission(&common.PermissionLevel{Actor: update.Account, Permission: update.Parent})
		minPermission = permission
	}
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckUpdateAuthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*minPermission, permissionIndex), &IrrelevantAuthException{},
		"UpdateAuth action declares irrelevant authority '%v'; minimum authority is %v", auth, common.PermissionLevel{update.Account, minPermission.Name})
}

func (a *AuthorizationManager) CheckDeleteAuthAuthorization(del DeleteAuth, auths []common.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "DeleteAuth action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == del.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")
	minPermission := a.GetPermission(&common.PermissionLevel{Actor: del.Account, Permission: del.Permission})
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckDeleteAuthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*minPermission, permissionIndex), &IrrelevantAuthException{},
		"DeleteAuth action declares irrelevant authority '%v'; minimum authority is %v", auth, common.PermissionLevel{minPermission.Owner, minPermission.Name})
}

func (a *AuthorizationManager) CheckLinkAuthAuthorization(link LinkAuth, auths []common.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "link action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == link.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")

	EosAssert(link.Type != UpdateAuth{}.GetName(), &ActionValidateException{}, "Cannot link eosio::updateauth to a minimum permission")
	EosAssert(link.Type != DeleteAuth{}.GetName(), &ActionValidateException{}, "Cannot link eosio::deleteauth to a minimum permission")
	EosAssert(link.Type != LinkAuth{}.GetName(), &ActionValidateException{}, "Cannot link eosio::linkauth to a minimum permission")
	EosAssert(link.Type != UnLinkAuth{}.GetName(), &ActionValidateException{}, "Cannot link eosio::unlinkauth to a minimum permission")
	EosAssert(link.Type != CancelDelay{}.GetName(), &ActionValidateException{}, "Cannot link eosio::canceldelay to a minimum permission")

	linkedPermissionName := a.LookupMinimumPermission(link.Account, link.Code, link.Type)
	if linkedPermissionName.Empty() {
		return
	}
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckLinkAuthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&common.PermissionLevel{link.Account, *linkedPermissionName}), permissionIndex), &IrrelevantAuthException{},
		"LinkAuth action declares irrelevant authority '%v'; minimum authority is %v", auth, common.PermissionLevel{link.Account, *linkedPermissionName})
}

func (a *AuthorizationManager) CheckUnLinkAuthAuthorization(unlink UnLinkAuth, auths []common.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "unlink action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == unlink.Account, &IrrelevantAuthException{},
		"the owner of the affected permission needs to be the actor of the declared authorization")

	unlinkedPermissionName := a.LookupLinkedPermission(unlink.Account, unlink.Code, unlink.Type)
	EosAssert(!unlinkedPermissionName.Empty(), &TransactionException{},
		"cannot unlink non-existent permission link of account '%v' for actions matching '%v::%v", unlink.Account, unlink.Code, unlink.Type)

	if *unlinkedPermissionName == common.DefaultConfig.EosioAnyName {
		return
	}
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckUnLinkAuthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&common.PermissionLevel{unlink.Account, *unlinkedPermissionName}), permissionIndex), &IrrelevantAuthException{},
		"unlink action declares irrelevant authority '%v'; minimum authority is %v", auth, common.PermissionLevel{unlink.Account, *unlinkedPermissionName})
}

func (a *AuthorizationManager) CheckCancelDelayAuthorization(cancel CancelDelay, auths []common.PermissionLevel) common.Microseconds {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "CancelDelay action should only have one declared authorization")
	auth := auths[0]
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckCancelDelayAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&cancel.CancelingAuth), permissionIndex), &IrrelevantAuthException{},
		"CancelDelay action declares irrelevant authority '%v'; specified authority to satisfy is %v", auth, cancel.CancelingAuth)

	generatedTrx := entity.GeneratedTransactionObject{}
	trxId := cancel.TrxId
	generatedIndex, err := a.control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
	if err != nil {
		log.Error("CheckCancelDelayAuthorization is error: %s", err)
	}
	itr, err := generatedIndex.LowerBound(entity.GeneratedTransactionObject{TrxId: trxId})
	if err != nil {
		log.Error("CheckCancelDelayAuthorization is error: %s", err)
	}

	err = itr.Data(&generatedTrx)
	EosAssert(err == nil && generatedTrx.TrxId == trxId, &TxNotFound{},
		"cannot cancel trx_id=%v, there is no deferred transaction with that transaction id", trxId)

	trx := types.Transaction{}
	rlp.DecodeBytes(generatedTrx.PackedTrx, &trx)
	found := false
	for _, act := range trx.Actions {
		for _, auth := range act.Authorization {
			if auth == cancel.CancelingAuth {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	EosAssert(found, &ActionValidateException{}, "canceling_auth in CancelDelay action was not found as authorization in the original delayed transaction")
	return common.Milliseconds(int64(generatedTrx.DelayUntil) - int64(generatedTrx.Published))
}

func (a *AuthorizationManager) CheckAuthorization(actions []*types.Action,
	providedKeys *PublicKeySet,
	providedPermissions *PermissionLevelSet,
	providedDelay common.Microseconds,
	checkTime *func(),
	allowUnusedKeys bool,
) {
	delayMaxLimit := common.Seconds(int64(a.control.GetGlobalProperties().Configuration.MaxTrxDelay))
	var effectiveProvidedDelay common.Microseconds
	if providedDelay >= delayMaxLimit {
		effectiveProvidedDelay = common.MaxMicroseconds()
	} else {
		effectiveProvidedDelay = providedDelay
	}
	checker := types.MakeAuthChecker(func(p *common.PermissionLevel) types.SharedAuthority {
		perm := a.GetPermission(p)
		if perm != nil {
			return perm.Auth
		} else {
			return types.SharedAuthority{}
		}
	},
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		providedKeys,
		providedPermissions,
		effectiveProvidedDelay,
		checkTime,
	)
	permissionToSatisfy := make(map[common.PermissionLevel]common.Microseconds)

	for _, act := range actions {
		specialCase := false
		delay := effectiveProvidedDelay

		if act.Account == common.DefaultConfig.SystemAccountName {
			specialCase = true
			switch act.Name {
			case UpdateAuth{}.GetName():
				UpdateAuth := UpdateAuth{}
				rlp.DecodeBytes(act.Data, &UpdateAuth)
				a.CheckUpdateAuthAuthorization(UpdateAuth, act.Authorization)

			case DeleteAuth{}.GetName():
				DeleteAuth := DeleteAuth{}
				rlp.DecodeBytes(act.Data, &DeleteAuth)
				a.CheckDeleteAuthAuthorization(DeleteAuth, act.Authorization)

			case LinkAuth{}.GetName():
				LinkAuth := LinkAuth{}
				rlp.DecodeBytes(act.Data, &LinkAuth)
				a.CheckLinkAuthAuthorization(LinkAuth, act.Authorization)

			case UnLinkAuth{}.GetName():
				UnLinkAuth := UnLinkAuth{}
				rlp.DecodeBytes(act.Data, &UnLinkAuth)
				a.CheckUnLinkAuthAuthorization(UnLinkAuth, act.Authorization)

			case CancelDelay{}.GetName():
				CancelDelay := CancelDelay{}
				rlp.DecodeBytes(act.Data, &CancelDelay)
				a.CheckCancelDelayAuthorization(CancelDelay, act.Authorization)

			default:
				specialCase = false
			}
		}

		for _, declaredAuth := range act.Authorization {
			(*checkTime)()
			if !specialCase {
				minPermissionName := a.LookupMinimumPermission(declaredAuth.Actor, act.Account, act.Name)
				if minPermissionName != nil {
					minPermission := a.GetPermission(&common.PermissionLevel{Actor: declaredAuth.Actor, Permission: *minPermissionName})
					permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
					if err != nil {
						log.Error("CheckAuthorization is error: %s", err)
					}
					EosAssert(a.GetPermission(&declaredAuth).Satisfies(*minPermission, permissionIndex), &IrrelevantAuthException{},
						"action declares irrelevant authority '%v'; minimum authority is %v", declaredAuth, common.PermissionLevel{minPermission.Owner, minPermission.Name})
				}
			}

			isExist := false
			for first, second := range permissionToSatisfy {
				if first == declaredAuth {
					if second > delay {
						second = delay
						isExist = true
						break
					}
				}
			}
			if !isExist {
				permissionToSatisfy[declaredAuth] = delay
			}
		}
	}
	for p, q := range permissionToSatisfy {
		(*checkTime)()
		EosAssert(checker.SatisfiedLoc(&p, q, nil), &UnsatisfiedAuthorization{},
			"transaction declares authority '%v', "+
				"but does not have signatures for it under a provided delay of %v ms, "+
				"provided permissions %v, and provided keys %v", p, providedDelay.Count()/1000, providedPermissions, providedKeys)
	}
	if !allowUnusedKeys {
		EosAssert(checker.AllKeysUsed(), &TxIrrelevantSig{}, "transaction bears irrelevant signatures from these keys: %v", checker.GetUnusedKeys())
	}
}

func (a *AuthorizationManager) CheckAuthorization2(account common.AccountName,
	permission common.PermissionName,
	providedKeys *PublicKeySet, //flat_set<public_key_type>
	providedPermissions *PermissionLevelSet, //flat_set<permission_level>
	providedDelay common.Microseconds,
	checkTime *func(),
	allowUnusedKeys bool,
) {
	delayMaxLimit := common.Seconds(int64(a.control.GetGlobalProperties().Configuration.MaxTrxDelay))
	var effectiveProvidedDelay common.Microseconds
	if providedDelay >= delayMaxLimit {
		effectiveProvidedDelay = common.MaxMicroseconds()
	} else {
		effectiveProvidedDelay = providedDelay
	}
	checker := types.MakeAuthChecker(func(p *common.PermissionLevel) types.SharedAuthority {
		perm := a.GetPermission(p)
		if perm != nil {
			return perm.Auth
		} else {
			return types.SharedAuthority{}
		}
	},
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		providedKeys,
		providedPermissions,
		effectiveProvidedDelay,
		checkTime,
	)
	EosAssert(checker.SatisfiedLc(&common.PermissionLevel{account, permission}, nil), &UnsatisfiedAuthorization{},
		"permission '%v' was not satisfied under a provided delay of %v ms, provided permissions %v, and provided keys %v",
		common.PermissionLevel{account, permission}, providedDelay.Count()/1000, providedPermissions, providedKeys)

	if !allowUnusedKeys {
		EosAssert(checker.AllKeysUsed(), &TxIrrelevantSig{}, "irrelevant keys provided: %v", checker.GetUnusedKeys())
	}
}

func (a *AuthorizationManager) GetRequiredKeys(trx *types.Transaction,
	candidateKeys *PublicKeySet,
	providedDelay common.Microseconds) PublicKeySet {
	checker := types.MakeAuthChecker(
		func(p *common.PermissionLevel) types.SharedAuthority {
			perm := a.GetPermission(p)
			if perm != nil {
				return perm.Auth
			} else {
				return types.SharedAuthority{}
			}
		},
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		candidateKeys,
		NewPermissionLevelSet(),
		providedDelay,
		noopCheckTime,
	)
	for _, act := range trx.Actions {
		for _, declaredAuth := range act.Authorization {
			EosAssert(checker.SatisfiedLc(&declaredAuth, nil), &UnsatisfiedAuthorization{},
				"transaction declares authority '%v', but does not have signatures for it.", declaredAuth)
		}
	}
	return checker.GetUsedKeys()
}
