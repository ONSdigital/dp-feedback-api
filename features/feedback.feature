Feature: Feedback


  Scenario: Posting valid feedback with all fields
    Given I am authorised
    When I POST "/feedback"
      """
        {
          "is_page_useful": true,
          "is_general_feedback": false,
          "ons_url": "https://localhost/subpath/one",
          "feedback": "very nice and useful website!",
          "name": "Mr Reporter",
          "email_address": "feedback@reporter.com"
        }
      """
    Then I should receive a 201 status code with an empty body response
    And The following email is sent
      """
        From: sender@feedback.com
        To: receiver@feedback.com
        Subject: Feedback received

        Feedback Type: A specific page
        Page URL: https://localhost/subpath/one
        Description: very nice and useful website!
        Name: Mr Reporter
        Email address: feedback@reporter.com
      """


  Scenario: Posting valid general feedback
    Given I am authorised
    When I POST "/feedback"
      """
        {
          "is_page_useful": false,
          "is_general_feedback": true,
          "feedback": "very nice and useful website!",
          "name": "Mr Reporter",
          "email_address": "feedback@reporter.com"
        }
      """
    Then I should receive a 201 status code with an empty body response
    And The following email is sent
      """
        From: sender@feedback.com
        To: receiver@feedback.com
        Subject: Feedback received

        Description: very nice and useful website!
        Name: Mr Reporter
        Email address: feedback@reporter.com
      """


  Scenario: Posting valid feedback with only required fields
    Given I am authorised
    When I POST "/feedback"
      """
        {
          "is_page_useful": true,
          "is_general_feedback":false,
          "feedback": "very nice and useful page!"
        }
      """
    Then I should receive a 201 status code with an empty body response
    And The following email is sent
      """
        From: sender@feedback.com
        To: receiver@feedback.com
        Subject: Feedback received

        Feedback Type: A specific page
        Description: very nice and useful page!
      """


  Scenario: Posting feedback with the wrong domain in ons url
    Given I am authorised
    When I POST "/feedback"
      """
        {
          "is_page_useful": false,
          "is_general_feedback": true,
          "ons_url": "https://attacker/subpath/one"
        }
      """
    Then I should receive a 400 status code with an the following body response
      """
        Key: 'Feedback.OnsURL' Error:Field validation for 'OnsURL' failed on the 'ons_url' tag
      """
    And No email is sent


  Scenario: Posting feedback with an invalid email address
    Given I am authorised
    When I POST "/feedback"
      """
        {
          "is_page_useful": false,
          "is_general_feedback": true,
          "email_address": "wrong.format"
        }
      """
    Then I should receive a 400 status code with an the following body response
      """
        Key: 'Feedback.EmailAddress' Error:Field validation for 'EmailAddress' failed on the 'email' tag
      """
    And No email is sent


  Scenario: Posting feedback with unsafe strings
    Given I am authorised
    When I POST "/feedback"
      """
        {
          "is_page_useful": true,
          "is_general_feedback": false,
          "ons_url": "https://localhost/subpath/one",
          "feedback": "<script>document.getElementById('demo').innerHTML = 'Hello JavaScript!'';</script>"
        }
      """
    Then I should receive a 201 status code with an empty body response
    And The following email is sent
      """
        From: sender@feedback.com
        To: receiver@feedback.com
        Subject: Feedback received

        Feedback Type: A specific page
        Page URL: https://localhost/subpath/one
        Description: &lt;script&gt;document.getElementById(\&#39;demo\&#39;).innerHTML = \&#39;Hello JavaScript!\&#39;\&#39;;&lt;/script&gt;
      """
