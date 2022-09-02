```
export TF_CLI_ARGS_init="-backend-config=dev.backend"
export TF_CLI_ARGS="-var-file=dev.tfvars"
```

`-i`, `--auto-init` Run `terraform init` automatically.

`-c`, `--config` Use config file.

`-q`, `--quite` Do not check if entry/file exists.

```
tfspace cd -i
tfspace backend add <env> <file>
tfspace backend rm <env> <file>
tfspace varfile add <env> <file>
tfspace varfile rm <env> <file>
tfspace edit # Edit config
tfspace validate # Validate config
tfspace env # Show current environments
```
