Feature: show help
  Terraformer needs to be able to access help

  Scenario: tfspace help
    When Terraformer runs "tfspace help"
    Then tfspace should print "help" content on screen

  Scenario: tfspace --help
    When Terraformer runs "tfspace --help"
    Then tfspace should print "help" content on screen

  Scenario: tfspace abc
    When Terraformer runs "tfspace abc"
    Then tfspace should print "space not found" error on screen
