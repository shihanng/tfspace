```
export TF_CLI_ARGS_init="-backend-config=dev.backend"
export TF_CLI_ARGS_plan="-var-file=dev.tfvars"
export TF_CLI_ARGS_apply="-var-file=dev.tfvars"
export TF_WORKSPACE=dev
```

`-i`, `--auto-init` Run `terraform init` automatically.

`-c`, `--config` Use config file.

`-n`, `--no-verify` Do not check if entry/file exists.

```
tfspace <env> -i
tfspace rm <env>
tfspace workspace add <env> <file>
tfspace workspace rm <env> <file>
tfspace backend add <env> <file>
tfspace backend rm <env> <file>
tfspace varfile add <env> <file>
tfspace varfile rm <env> <file> # Remove if everything is rm
tfspace edit # Edit config
tfspace validate # Validate config
tfspace env # Show current environments

TFSPACE=<ENV>
```

- Auto completion?
- Shell prompt?
