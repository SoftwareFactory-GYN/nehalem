# Nehalem

> Code quality metrics run in isolated docker containers.

This software allows engineers and managers to assert the quality of their code. 
This is all done in a distributed, high performance, isolated manner.

All builds are run inside docker containers, and build can be triggered manually or automatically via 
git web hooks, metrics must be echoed out of the build as json.  

Currently we are working with two kinds of metrics, code coverage, which is a measure used to describe the 
degree to which the source code of a program is executed when a particular test suite runs, and lint, which 
refers to tools that analyze source code to flag programming errors, bugs, stylistic errors, and suspicious 
constructs. With these metrics in hand we can supply graphs to aid and support decision making and code 
quality assurance.



Current supported languages include:
- [x] Python
- [x] GoLang

With support coming for the following languages:
- [ ] Javascript
- [ ] Java
