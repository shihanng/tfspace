Feature: workspace
  Terraformer needs to be able to add or remove workspace

  Scenario: add workspace
    Given a project without tfspace.yml
    When Terraformer runs "tfspace workspace add dev development"
    And Terraformer runs "tfspace workspace add stg staging"
    Then the tfspace.yml should contain:
      """
      dev:
        workspace: development
      stg:
        workspace: staging
      """
