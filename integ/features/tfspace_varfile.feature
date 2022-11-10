Feature: varfile
  Terraformer needs to be able to add or remove varfile

  Scenario: add or remove workspace
    Given a project without tfspace.yml
    When Terraformer runs "tfspace varfile add dev dev.tfvars"
    Then tfspace should run without error
    When Terraformer runs "tfspace varfile add stg stg.tfvars"
    Then tfspace should run without error
    And the tfspace.yml should contain:
      """
      dev:
        varfile:
          - dev.tfvars
      stg:
        varfile:
          - stg.tfvars
    
      """
    When Terraformer runs "tfspace varfile add dev development.tfvars"
    Then the tfspace.yml should contain:
      """
      dev:
        varfile:
          - dev.tfvars
          - development.tfvars
      stg:
        varfile:
          - stg.tfvars

      """
    When Terraformer runs "tfspace varfile rm dev dev.tfvars"
    And Terraformer runs "tfspace varfile rm stg stg.tfvars"
    And Terraformer runs "tfspace varfile rm prod prod.tfvars"
    Then tfspace should run without error
    And the tfspace.yml should contain:
      """
      dev:
        varfile:
          - development.tfvars

      """
    When Terraformer runs "tfspace use dev -a" and then env
    Then tfspace should run without error
    And should set environment variables:
      | TFSPACE           | dev                            |
      | TF_CLI_ARGS_plan  | '-var-file=development.tfvars' |
      | TF_CLI_ARGS_apply | '-var-file=development.tfvars' |
