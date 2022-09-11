Feature: workspace
  Terraformer needs to be able to add or remove workspace

  Scenario: add workspace
    Given a project without tfspace.yml
    When Terraformer runs "tfspace workspace add dev development"
    Then tfspace should run without error
    And the tfspace.yml should contain:
      """
      dev:
        backend: []
        varfile: []
        workspace: development

      """
