package cmd

const (
	defaultMgmtDataDir   = "/var/lib/cosmos/"
	defaultMgmtConfigDir = "/etc/cosmos"
	defaultLogDir        = "/var/log/cosmos"

	oldDefaultMgmtDataDir   = "/var/lib/wiretrustee/"
	oldDefaultMgmtConfigDir = "/etc/wiretrustee"
	oldDefaultLogDir        = "/var/log/wiretrustee"

	defaultMgmtConfig    = defaultMgmtConfigDir + "/management.json"
	defaultLogFile       = defaultLogDir + "/management.log"
	oldDefaultMgmtConfig = oldDefaultMgmtConfigDir + "/management.json"
	oldDefaultLogFile    = oldDefaultLogDir + "/management.log"

	defaultSingleAccModeDomain = "cosmos.selfhosted"
)
