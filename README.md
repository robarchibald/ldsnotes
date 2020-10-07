# ldsnotes
[![Build Status](https://travis-ci.org/robarchibald/ldsnotes.svg?branch=master)](https://travis-ci.org/robarchibald/ldsnotes) [![Coverage Status](https://coveralls.io/repos/github/robarchibald/ldsnotes/badge.svg?branch=master)](https://coveralls.io/github/robarchibald/ldsnotes?branch=master)

A Go tool to parse and format the individual notes that a person can save on churchofjesuschrist.org.

** Setup
1. Restore the scriptures database to your local Postgres database, or find your database of choice at https://scriptures.nephi.org/
2. `psql -d scriptures -f lds-scriptures-psql.sql`

** Getting json file for parsing

1. Go to churchofjesuschrist.org
2. Click on "My Account and Ward" in the upper right hand corner of the page
3. Sign in to your account
4. From the "My Account and Ward" menu, select "Notes"
  - Open developer tools and find this API call: https://www.churchofjesuschrist.org/notes/api/v2
  - Paste it into a separate browser window
  - Change the number to return to something higher to get all records at once
  - Copy the text into a json file
5. OR go to https://www.churchofjesuschrist.org/notes/api/v2/annotations?notesAsHtml=true&numberToReturn=5000&tags=REPLACE_WITH_TAG&type=highlight%2Cjournal%2Creference replacing REPLACE_WITH_TAG with the desired tag 