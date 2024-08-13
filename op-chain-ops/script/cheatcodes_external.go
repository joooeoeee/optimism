package script

import (
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum-optimism/optimism/op-chain-ops/foundry"
	"github.com/ethereum/go-ethereum/core/vm"
)

// Ffi implements https://book.getfoundry.sh/cheatcodes/ffi
func (c *CheatCodesPrecompile) Ffi(args []string) ([]byte, error) {
	return nil, vm.ErrExecutionReverted
}

// Prompt implements https://book.getfoundry.sh/cheatcodes/prompt
func (c *CheatCodesPrecompile) Prompt() error {
	return vm.ErrExecutionReverted
}

// ProjectRoot implements https://book.getfoundry.sh/cheatcodes/project-root
func (c *CheatCodesPrecompile) ProjectRoot() string {
	return ""
}

func (c *CheatCodesPrecompile) getArtifact(input string) (*foundry.Artifact, error) {
	// TODO fetching by relative file path, or using a contract version, is not supported
	parts := strings.SplitN(input, ":", 1)
	name := parts[0] + ".sol"
	contract := parts[0]
	if len(parts) == 2 {
		name = parts[0]
		contract = parts[1]
	}
	return c.h.af.ReadArtifact(name, contract)
}

// GetCode implements https://book.getfoundry.sh/cheatcodes/get-code
func (c *CheatCodesPrecompile) GetCode(input string) ([]byte, error) {
	artifact, err := c.getArtifact(input)
	if err != nil {
		return nil, err
	}
	return artifact.Bytecode.Object, nil
}

// GetDeployedCode implements https://book.getfoundry.sh/cheatcodes/get-deployed-code
func (c *CheatCodesPrecompile) GetDeployedCode(input string) ([]byte, error) {
	artifact, err := c.getArtifact(input)
	if err != nil {
		return nil, err
	}
	return artifact.DeployedBytecode.Object, nil
}

// Sleep implements https://book.getfoundry.sh/cheatcodes/sleep
func (c *CheatCodesPrecompile) Sleep(ms *big.Int) error {
	if !ms.IsUint64() {
		return vm.ErrExecutionReverted
	}
	time.Sleep(time.Duration(ms.Uint64()) * time.Millisecond)
	return nil
}

// UnixTime implements https://book.getfoundry.sh/cheatcodes/unix-time
func (c *CheatCodesPrecompile) UnixTime() (ms *big.Int) {
	return big.NewInt(time.Now().UnixMilli())
}

// SetEnv implements https://book.getfoundry.sh/cheatcodes/set-env
func (c *CheatCodesPrecompile) SetEnv(key string, value string) error {
	if key == "" {
		return errors.New("env key must not be empty")
	}
	if strings.ContainsRune(key, '=') {
		return errors.New("env key must not contain = sign")
	}
	if strings.ContainsRune(key, 0) {
		return errors.New("env key must not contain NUL")
	}
	if strings.ContainsRune(value, 0) {
		return errors.New("env value must not contain NUL")
	}
	c.h.envVars[key] = value
	return nil
}

// EnvOr implements https://book.getfoundry.sh/cheatcodes/env-or
func (c *CheatCodesPrecompile) EnvOr() {
	// TODO
}

// EnvBool implements https://book.getfoundry.sh/cheatcodes/env-bool
func (c *CheatCodesPrecompile) EnvBool() {

}

// EnvUint implements https://book.getfoundry.sh/cheatcodes/env-uint
func (c *CheatCodesPrecompile) EnvUint() {

}

// EnvInt implements https://book.getfoundry.sh/cheatcodes/env-int
func (c *CheatCodesPrecompile) EnvInt() {

}

// EnvAddress implements https://book.getfoundry.sh/cheatcodes/env-address
func (c *CheatCodesPrecompile) EnvAddress() {

}

// EnvBytes32 implements https://book.getfoundry.sh/cheatcodes/env-bytes32
func (c *CheatCodesPrecompile) EnvBytes32() {

}

// EnvString implements https://book.getfoundry.sh/cheatcodes/env-string
func (c *CheatCodesPrecompile) EnvString() {

}

// EnvBytes implements https://book.getfoundry.sh/cheatcodes/env-bytes
func (c *CheatCodesPrecompile) EnvBytes() {

}

// KeyExists implements https://book.getfoundry.sh/cheatcodes/key-exists
func (c *CheatCodesPrecompile) KeyExists(jsonData string, key string) (bool, error) {
	return c.KeyExistsJson(jsonData, key)

}

// KeyExistsJson implements https://book.getfoundry.sh/cheatcodes/key-exists-json
func (c *CheatCodesPrecompile) KeyExistsJson(jsonData string, key string) (bool, error) {
	var x map[string]json.RawMessage
	if err := json.Unmarshal([]byte(jsonData), &x); err != nil {
		return false, err
	}
	_, ok := x[key]
	return ok, nil
}

// KeyExistsToml implements https://book.getfoundry.sh/cheatcodes/key-exists-toml
func (c *CheatCodesPrecompile) KeyExistsToml() {
	// TODO
}

// ParseJSON implements https://book.getfoundry.sh/cheatcodes/parse-json
func (c *CheatCodesPrecompile) ParseJson() {
	// TODO
}

// ParseToml implements https://book.getfoundry.sh/cheatcodes/parse-toml
func (c *CheatCodesPrecompile) ParseToml() {
	// TODO
}

// ParseJsonKeys implements https://book.getfoundry.sh/cheatcodes/parse-json-keys
func (c *CheatCodesPrecompile) ParseJsonKeys(jsonData string, key string) []string {
	// TODO jq like behavior wih key
	return nil
}

// ParseTomlKeys implements https://book.getfoundry.sh/cheatcodes/parse-toml-keys
func (c *CheatCodesPrecompile) ParseTomlKeys() {
	// TODO
}

// SerializeJson implements https://book.getfoundry.sh/cheatcodes/serialize-json
func (c *CheatCodesPrecompile) SerializeJson() {
	// TODO
}

// WriteJson implements https://book.getfoundry.sh/cheatcodes/write-json
func (c *CheatCodesPrecompile) WriteJson() {
	// TODO
}

// WriteToml implements https://book.getfoundry.sh/cheatcodes/write-toml
func (c *CheatCodesPrecompile) WriteToml() {
	// TODO
}
