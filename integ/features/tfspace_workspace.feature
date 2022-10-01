Feature: workspace
  Terraformer needs to be able to add or remove workspace

  Scenario: add or remove workspace
    Given a project without tfspace.yml
    When Terraformer runs "tfspace workspace add dev development"
    Then tfspace should run without error
    When Terraformer runs "tfspace workspace add stg staging"
    Then tfspace should run without error
    And the tfspace.yml should contain:
      """
      dev:
        workspace: development
      stg:
        workspace: staging

      """
    When Terraformer runs "tfspace workspace add dev dev"
    Then the tfspace.yml should contain:
      """
      dev:
        workspace: dev
      stg:
        workspace: staging

      """
    When Terraformer runs "tfspace workspace rm dev"
    And Terraformer runs "tfspace workspace rm dev"
    Then tfspace should run without error
    And the tfspace.yml should contain:
      """
      stg:
        workspace: staging

      """
