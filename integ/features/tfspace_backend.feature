Feature: backend
  Terraformer needs to be able to add or remove backend

  Scenario: add or remove workspace
    Given a project without tfspace.yml
    When Terraformer runs "tfspace backend add dev dev.backend"
    Then tfspace should run without error
    When Terraformer runs "tfspace backend add stg stg.backend"
    Then tfspace should run without error
    And the tfspace.yml should contain:
      """
      dev:
        backend:
          - dev.backend
      stg:
        backend:
          - stg.backend
    
      """
    When Terraformer runs "tfspace backend add dev development.backend"
    Then the tfspace.yml should contain:
      """
      dev:
        backend:
          - dev.backend
          - development.backend
      stg:
        backend:
          - stg.backend

      """
    When Terraformer runs "tfspace backend rm dev dev.backend"
    And Terraformer runs "tfspace backend rm stg stg.backend"
    Then tfspace should run without error
    And the tfspace.yml should contain:
      """
      dev:
        backend:
          - development.backend

      """
    When Terraformer runs "tfspace dev" and then env
    Then tfspace should run without error
    And should set environment variables:
      | TFSPACE          | dev                                     |
      | TF_CLI_ARGS_init | '-backend-config="development.backend"' |
