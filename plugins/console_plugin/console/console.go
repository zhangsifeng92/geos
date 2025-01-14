package console

import (
	"fmt"
	"github.com/peterh/liner"
	"github.com/robertkrimen/otto"
	"github.com/zhangsifeng92/geos/plugins/console_plugin/console/jsre"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

var (
	onlyWhitespace = regexp.MustCompile(`^\s*$`)
	exit           = regexp.MustCompile(`^\s*exit\s*;*\s*$`)
)

// HistoryFile is the file within the data directory to store input scrollback.
const HistoryFile = "history"

// DefaultPrompt is the default prompt line prefix to use for user input querying.
const DefaultPrompt = "eosgo > "

// Config is the collection of configurations to fine tune the behavior of the
// JavaScript console.
type Config struct {
	DataDir  string       // Data directory to store the console history at
	DocRoot  string       // Filesystem path from where to load JavaScript files from
	Client   *client      // RPC client to execute EOSGO requests through
	Prompt   string       // Input prompt prefix string (defaults to DefaultPrompt)
	Prompter UserPrompter // Input prompter to allow interactive user feedback (defaults to TerminalPrompter)
	Printer  io.Writer    // Output writer to serialize any display strings to (defaults to os.Stdout)
	Preload  []string     // Absolute paths to JavaScript files to preload
}

// Console is a JavaScript interpreted runtime environment. It is a fully fledged
// JavaScript console attached to a running node via an external or in-process RPC
// client.
type Console struct {
	client   *client      // RPC client to execute EOSGO requests through
	jsre     *jsre.JSRE   // JavaScript runtime environment running the interpreter
	prompt   string       // Input prompt prefix string
	prompter UserPrompter // Input prompter to allow interactive user feedback
	histPath string       // Absolute path to the console scrollback history
	history  []string     // Scroll history maintained by the console
	printer  io.Writer    // Output writer to serialize any display strings to
}

// New initializes a JavaScript interpreted runtime environment and sets defaults
// with the config struct.
func New(config Config) (*Console, error) {
	// Handle unset config values gracefully
	if config.Prompter == nil {
		config.Prompter = Stdin
	}
	if config.Prompt == "" {
		config.Prompt = DefaultPrompt
	}
	if config.Printer == nil {
		config.Printer = os.Stdout
	}
	// Initialize the console and return
	console := &Console{
		//client:   config.Client,
		client:   config.Client,
		jsre:     jsre.New(config.DocRoot, config.Printer),
		prompt:   config.Prompt,
		prompter: config.Prompter,
		printer:  config.Printer,
		histPath: filepath.Join(config.DataDir, HistoryFile),
	}
	if err := os.MkdirAll(config.DataDir, 0700); err != nil {
		return nil, err
	}

	if err := console.init(config.Preload); err != nil {
		return nil, err
	}

	return console, nil
}

// init retrieves the available APIs from the remote RPC provider and initializes
// the console's JavaScript namespaces based on the exposed modules.
func (c *Console) init(preload []string) (err error) {
	// Initialize the JavaScript <-> Go RPC bridge
	consoleObj, _ := c.jsre.Get("console")
	consoleObj.Object().Set("log", c.consoleOutput)
	consoleObj.Object().Set("error", c.consoleOutput)

	eos := newEosgo(c)
	//bridge := newBridge(c.client, c.prompter, c.printer)
	c.jsre.Bind("eos", eos)

	chain := newchainAPI(c)
	c.jsre.Bind("chain", chain)
	//c.jsre.Set("chain", struct{}{})
	//chainObj, _ := c.jsre.Get("chain")
	//chainObj.Object().Set("getInfo", chain.GetInfo)
	//chainObj.Object().Set("getBlock", chain.GetBlock)
	//chainObj.Object().Set("getAccount", chain.GetAccount)
	//chainObj.Object().Set("getCode", chain.GetCode)
	//chainObj.Object().Set("getAbi", chain.GetAbi)

	wallet := newWalletApi(c)
	c.jsre.Bind("wallet", wallet)
	//c.jsre.Set("createKey", eos.CreateKey)
	//c.jsre.Set("wallet", struct{}{})
	//walletObj, _ := c.jsre.Get("wallet")
	//walletObj.Object().Set("createWallet", wallet.CreateWallet)
	//walletObj.Object().Set("importKey",wallet.ImportKey)

	net := newNetAPI(c)
	c.jsre.Bind("net", net)

	system := newSystem(c)
	c.jsre.Bind("system", system)

	producer := newProduceAPI(c)
	c.jsre.Bind("producer", producer)

	multiSig := newMultiSig(c)
	//c.jsre.Bind("multiSig", multiSig)
	c.jsre.Set("multisig", struct{}{})
	multiSigObj, _ := c.jsre.Get("multisig")
	multiSigObj.Object().Set("propose", multiSig.Propose)
	multiSigObj.Object().Set("proposetrx", multiSig.ProposeTrx)
	multiSigObj.Object().Set("review", multiSig.Review)
	multiSigObj.Object().Set("approve", multiSig.Approve)
	multiSigObj.Object().Set("unapprove", multiSig.Unapprove)
	multiSigObj.Object().Set("cancel", multiSig.Cancel)
	multiSigObj.Object().Set("exec", multiSig.Exec)

	//if err := c.jsre.Compile("eosgo.js", jsre.Eosgo_JS); err != nil {
	//	return fmt.Errorf("eosgo.js:%v", err)
	//}
	//

	//The admin.sleep and admin.sleepBlocks are offered by the console and not by the RPC layer.
	admin, err := c.jsre.Get("eos")
	if err != nil {
		return err
	}
	if obj := admin.Object(); obj != nil { // make sure the admin api is enabled over the interface
		obj.Set("clearHistory", c.clearHistory)
	}
	// Preload any JavaScript files before starting the console
	for _, path := range preload {
		if err := c.jsre.Exec(path); err != nil {
			failure := err.Error()
			if ottoErr, ok := err.(*otto.Error); ok {
				failure = ottoErr.String()
			}
			return fmt.Errorf("%s: %v", path, failure)
		}
	}
	// Configure the console's input prompter for scrollback and tab completion
	if c.prompter != nil {
		if content, err := ioutil.ReadFile(c.histPath); err != nil {
			c.prompter.SetHistory(nil)
		} else {
			c.history = strings.Split(string(content), "\n")
			c.prompter.SetHistory(c.history)
		}
		c.prompter.SetWordCompleter(c.AutoCompleteInput)
	}
	return nil
}

func (c *Console) clearHistory() {
	c.history = nil
	c.prompter.ClearHistory()
	if err := os.Remove(c.histPath); err != nil {
		fmt.Fprintln(c.printer, "can't delete history file:", err)
	} else {
		fmt.Fprintln(c.printer, "history file deleted")
	}
}

// consoleOutput is an override for the console.log and console.error methods to
// stream the output into the configured output stream instead of stdout.
func (c *Console) consoleOutput(call otto.FunctionCall) otto.Value {
	output := []string{}
	for _, argument := range call.ArgumentList {
		output = append(output, fmt.Sprintf("%v", argument))
	}
	fmt.Fprintln(c.printer, strings.Join(output, " "))
	return otto.Value{}
}

// AutoCompleteInput is a pre-assembled word completer to be used by the user
// input prompter to provide hints to the user about the methods available.
func (c *Console) AutoCompleteInput(line string, pos int) (string, []string, string) {
	// No completions can be provided for empty inputs
	if len(line) == 0 || pos == 0 {
		return "", nil, ""
	}
	// Chunck data to relevant part for autocompletion
	// E.g. in case of nested lines eth.getBalance(eth.coinb<tab><tab>
	start := pos - 1
	for ; start > 0; start-- {
		// Skip all methods and namespaces (i.e. including the dot)
		if line[start] == '.' || (line[start] >= 'a' && line[start] <= 'z') || (line[start] >= 'A' && line[start] <= 'Z') {
			continue
		}
		// Handle web3 in a special way (i.e. other numbers aren't auto completed)
		if start >= 3 && line[start-3:start] == "web3" {
			start -= 3
			continue
		}
		// We've hit an unexpected character, autocomplete form here
		start++
		break
	}
	return line[:start], c.jsre.CompleteKeywords(line[start:pos]), line[pos:]
}

// Welcome show summary of current eosgo instance and some metadata about the
// console's available modules.
func (c *Console) Welcome() {
	// Print some generic eosgo metadata
	c.logo()

	//c.jsre.Run(`
	//        console.log("get info: " + chain.getInfo());
	//`)
}

// Interactive starts an interactive user session, where input is propted from
// the configured user prompter.
func (c *Console) Interactive() {
	var (
		prompt    = c.prompt          // Current prompt line (used for multi-line inputs)
		indents   = 0                 // Current number of input indents (used for multi-line inputs)
		input     = ""                // Current user input
		scheduler = make(chan string) // Channel to send the next prompt on and receive the input
	)
	// Start a goroutine to listen for promt requests and send back inputs
	go func() {
		for {
			// Read the next user input
			line, err := c.prompter.PromptInput(<-scheduler)
			if err != nil {
				// In case of an error, either clear the prompt or fail
				if err == liner.ErrPromptAborted { // ctrl-C
					prompt, indents, input = c.prompt, 0, ""
					scheduler <- ""
					continue
				}
				close(scheduler)
				return
			}
			// User input retrieved, send for interpretation and loop
			scheduler <- line
		}
	}()
	// Monitor Ctrl-C too in case the input is empty and we need to bail
	abort := make(chan os.Signal, 1)
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// Start sending prompts to the user and reading back inputs
		for {
			// Send the next prompt, triggering an input read and process the result
			scheduler <- prompt
			select {
			case <-abort:
				// User forcefully quite the console
				fmt.Fprintln(c.printer, "caught interrupt, exiting")
				return

			case line, ok := <-scheduler:
				// User input was returned by the prompter, handle special cases
				if !ok || (indents <= 0 && exit.MatchString(line)) {
					return
				}
				if onlyWhitespace.MatchString(line) {
					continue
				}
				// Append the line to the input and check for multi-line interpretation
				input += line + "\n"

				indents = countIndents(input)
				if indents <= 0 {
					prompt = c.prompt
				} else {
					prompt = strings.Repeat(" ", indents*3) + " "
				}
				// If all the needed lines are present, save the command and run
				if indents <= 0 {
					if len(input) > 0 && input[0] != ' ' {
						if command := strings.TrimSpace(input); len(c.history) == 0 || command != c.history[len(c.history)-1] {
							c.history = append(c.history, command)
							if c.prompter != nil {
								c.prompter.AppendHistory(command)
							}
						}
					}
					//fmt.Println("*************console input:  ", input)
					c.Evaluate(input)
					input = ""
				}
			}
		}
	}()
}

// Evaluate executes code and pretty prints the result to the specified output
// stream.
func (c *Console) Evaluate(statement string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(c.printer, "[native] error: %v\n", r)
		}
	}()
	return c.jsre.Evaluate(statement, c.printer)
}

// countIndents returns the number of identations for the given input.
// In case of invalid input such as var a = } the result can be negative.
func countIndents(input string) int {
	var (
		indents     = 0
		inString    = false
		strOpenChar = ' '   // keep track of the string open char to allow var str = "I'm ....";
		charEscaped = false // keep track if the previous char was the '\' char, allow var str = "abc\"def";
	)

	for _, c := range input {
		switch c {
		case '\\':
			// indicate next char as escaped when in string and previous char isn't escaping this backslash
			if !charEscaped && inString {
				charEscaped = true
			}
		case '\'', '"':
			if inString && !charEscaped && strOpenChar == c { // end string
				inString = false
			} else if !inString && !charEscaped { // begin string
				inString = true
				strOpenChar = c
			}
			charEscaped = false
		case '{', '(':
			if !inString { // ignore brackets when in string, allow var str = "a{"; without indenting
				indents++
			}
			charEscaped = false
		case '}', ')':
			if !inString {
				indents--
			}
			charEscaped = false
		default:
			charEscaped = false
		}
	}

	return indents
}

// Execute runs the JavaScript file specified as the argument.
func (c *Console) Execute(path string) error {
	return c.jsre.Exec(path)
}

// Stop cleans up the console and terminates the runtime environment.
func (c *Console) Stop(graceful bool) error {
	if err := ioutil.WriteFile(c.histPath, []byte(strings.Join(c.history, "\n")), 0600); err != nil {
		return err
	}
	if err := os.Chmod(c.histPath, 0600); err != nil { // Force 0600, even if it was different previously
		return err
	}
	c.jsre.Stop(graceful)
	return nil
}

func (c *Console) logo() {
	fmt.Fprintf(c.printer, "\x1b[1;31m\n")
	fmt.Fprintf(c.printer, "\t _______  _______  _______  _______  _______\n")
	fmt.Fprintf(c.printer, "\t(  ____ \\(  ___  )(  ____ \\(  ____ \\(  ___  )\n")
	fmt.Fprintf(c.printer, "\t| (    \\/| (   ) || (    \\/| (    \\/| (   ) |\n")
	fmt.Fprintf(c.printer, "\t| (__    | |   | || (_____ | |  ___ | |   | |\n")
	fmt.Fprintf(c.printer, "\t|  __)   | |   | |(_____  )| | (_  )| |   | |\n")
	fmt.Fprintf(c.printer, "\t| (      | |   | |      ) || |   ) || |   | |\n")
	fmt.Fprintf(c.printer, "\t| (____/\\| (___) |/\\____) || (___) || (___) |\n")
	fmt.Fprintf(c.printer, "\t(_______/(_______)\\_______)\\_______/(_______)\n")
	fmt.Fprintf(c.printer, "\x1b[0m\n")

	//fmt.Fprintf(c.printer, "\tFor more information:\n")
	//fmt.Fprintf(c.printer,"\tEOSGO Website: https://eos.io\n")
	//fmt.Fprintf(c.printer,"\tEOSGO Telegram channel @ https://t.me/EOSProject\n")
	//fmt.Fprintf(c.printer,"\tEOSGO Resources: https://eos.io/resources/\n")
	//fmt.Fprintf(c.printer,"\tEOSGO Stack Exchange: https://eosio.stackexchange.com\n")
	//fmt.Fprintf(c.printer, "\tGithub: https://github.com/zhangsifeng92/geos\n\n")
	fmt.Fprintf(c.printer, "\tWelcome to the EOSGO JavaScript console!\n")
}

// throwJSException panics on an otto.Value. The Otto VM will recover from the
// Go panic and throw msg as a JavaScript error.
func throwJSException(msg interface{}) otto.Value {
	val, err := otto.ToValue(msg)
	if err != nil {
		//log.Error("Failed to serialize JavaScript exception", "exception", msg, "err", err)
	}
	panic(val)
}
