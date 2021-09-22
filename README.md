# dispar

Dispar it's mean "Distributed Parser". 

The main idea is to get data very quickly using cloud capabilities. You just need to describe what data from which pages you want to get.

Add a declarative description of the range of page addresses given from these pages using html selectors. 

Then the cloud function is launched that collects the URLs of the pages for parsing and sends them to Google Pub / Sub. 

The received addresses are processed by the cloud function and sen data to the message queue. 

The GCE instance pulls them out and saves them to the cloud storage or database.
