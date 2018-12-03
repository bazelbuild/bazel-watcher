# Change Log

## [v0.8.2](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.8.2) (2018-12-03)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.8.1...v0.8.2)

**Merged pull requests:**

- Updates WORKSPACE file to work with bazel 0.20.0 [\#188](https://github.com/bazelbuild/bazel-watcher/pull/188) ([Jdban](https://github.com/Jdban))
- Print the list of supported Bazel flags [\#187](https://github.com/bazelbuild/bazel-watcher/pull/187) ([achew22](https://github.com/achew22))
- Update org\_golang\_x\_sys commit hash to 4ed8d59 [\#174](https://github.com/bazelbuild/bazel-watcher/pull/174) ([renovate[bot]](https://github.com/apps/renovate))

## [v0.8.1](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.8.1) (2018-12-02)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.8.0...v0.8.1)

**Implemented enhancements:**

- Too many files error [\#109](https://github.com/bazelbuild/bazel-watcher/issues/109)

**Closed issues:**

- Correct order of stdout/stdin [\#136](https://github.com/bazelbuild/bazel-watcher/issues/136)
- Deceptive warning message when running iBazel from npm dist [\#106](https://github.com/bazelbuild/bazel-watcher/issues/106)
- query for tags doesn't read through alias [\#100](https://github.com/bazelbuild/bazel-watcher/issues/100)
- Likely leaking file descriptors [\#96](https://github.com/bazelbuild/bazel-watcher/issues/96)
- installs error [\#90](https://github.com/bazelbuild/bazel-watcher/issues/90)
- ibazel not exiting on Ctrl-C on Mac OS after shutting down child process [\#50](https://github.com/bazelbuild/bazel-watcher/issues/50)
- publish gh release binaries [\#15](https://github.com/bazelbuild/bazel-watcher/issues/15)

**Merged pull requests:**

- Create a top level release dir [\#186](https://github.com/bazelbuild/bazel-watcher/pull/186) ([achew22](https://github.com/achew22))
- Fix npm warning logic [\#185](https://github.com/bazelbuild/bazel-watcher/pull/185) ([dougkoch](https://github.com/dougkoch))
- Update dependency io\_bazel\_rules\_go to v0.16.3 [\#183](https://github.com/bazelbuild/bazel-watcher/pull/183) ([renovate[bot]](https://github.com/apps/renovate))
- Extract an interface for fsnotify [\#182](https://github.com/bazelbuild/bazel-watcher/pull/182) ([achew22](https://github.com/achew22))
- Clean up tests a little [\#181](https://github.com/bazelbuild/bazel-watcher/pull/181) ([achew22](https://github.com/achew22))
- Uncross stdout and stderr buffers [\#138](https://github.com/bazelbuild/bazel-watcher/pull/138) ([dougkoch](https://github.com/dougkoch))

## [v0.8.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.8.0) (2018-11-28)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.7.0...v0.8.0)

**Closed issues:**

- Second and further file changes not detected on MacOS Mojave [\#173](https://github.com/bazelbuild/bazel-watcher/issues/173)
- Failure to install on Mac OS X Mojave via Homebrew [\#162](https://github.com/bazelbuild/bazel-watcher/issues/162)

**Merged pull requests:**

- Watch files via parent directory [\#180](https://github.com/bazelbuild/bazel-watcher/pull/180) ([dougkoch](https://github.com/dougkoch))
- Increase e2e timeout [\#179](https://github.com/bazelbuild/bazel-watcher/pull/179) ([dougkoch](https://github.com/dougkoch))
- Support Bazel --override\_repository flag [\#178](https://github.com/bazelbuild/bazel-watcher/pull/178) ([justbuchanan](https://github.com/justbuchanan))
- Use FSNotify's release tags [\#177](https://github.com/bazelbuild/bazel-watcher/pull/177) ([achew22](https://github.com/achew22))
- Push tags by address instead of by name [\#176](https://github.com/bazelbuild/bazel-watcher/pull/176) ([achew22](https://github.com/achew22))
- Support Bazel --define flag [\#175](https://github.com/bazelbuild/bazel-watcher/pull/175) ([IgorMinar](https://github.com/IgorMinar))
- Switch to tags in WORKSPACE [\#171](https://github.com/bazelbuild/bazel-watcher/pull/171) ([achew22](https://github.com/achew22))

## [v0.7.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.7.0) (2018-11-15)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.6.0...v0.7.0)

**Closed issues:**

- "Error getting Bazel info" on server start [\#146](https://github.com/bazelbuild/bazel-watcher/issues/146)
- ibazel watch stops working after 2 refreshes [\#117](https://github.com/bazelbuild/bazel-watcher/issues/117)
- Watches 0 files when target is in a different workspace [\#97](https://github.com/bazelbuild/bazel-watcher/issues/97)

**Merged pull requests:**

- Add release script to make releasing easy [\#172](https://github.com/bazelbuild/bazel-watcher/pull/172) ([achew22](https://github.com/achew22))
- Update org\_golang\_x\_sys commit hash to 66b7b13 [\#169](https://github.com/bazelbuild/bazel-watcher/pull/169) ([renovate[bot]](https://github.com/apps/renovate))
- Update com\_github\_gorilla\_websocket commit hash to 483fb8d [\#168](https://github.com/bazelbuild/bazel-watcher/pull/168) ([renovate[bot]](https://github.com/apps/renovate))
- Update com\_github\_golang\_protobuf commit hash to 951a149 [\#167](https://github.com/bazelbuild/bazel-watcher/pull/167) ([renovate[bot]](https://github.com/apps/renovate))
- Update com\_github\_fsnotify\_fsnotify commit hash to ccc981b [\#166](https://github.com/bazelbuild/bazel-watcher/pull/166) ([renovate[bot]](https://github.com/apps/renovate))
- Update com\_github\_bazelbuild\_rules\_go commit hash to 109c520 [\#165](https://github.com/bazelbuild/bazel-watcher/pull/165) ([renovate[bot]](https://github.com/apps/renovate))
- Warn if no files are being watched [\#164](https://github.com/bazelbuild/bazel-watcher/pull/164) ([dougkoch](https://github.com/dougkoch))
- Support Bazel --strategy flag [\#161](https://github.com/bazelbuild/bazel-watcher/pull/161) ([Globegitter](https://github.com/Globegitter))
- Update dependency io\_bazel\_rules\_go to v0.16.2 [\#160](https://github.com/bazelbuild/bazel-watcher/pull/160) ([renovate[bot]](https://github.com/apps/renovate))
- Invoke python explicitly in --workplace\_status\_command [\#159](https://github.com/bazelbuild/bazel-watcher/pull/159) ([achew22](https://github.com/achew22))
- If in windows, script\_path should end in .bat [\#158](https://github.com/bazelbuild/bazel-watcher/pull/158) ([achew22](https://github.com/achew22))
- Ignore startup message when parsing `bazel info` output [\#157](https://github.com/bazelbuild/bazel-watcher/pull/157) ([dougkoch](https://github.com/dougkoch))
- Fix README typo [\#156](https://github.com/bazelbuild/bazel-watcher/pull/156) ([dougkoch](https://github.com/dougkoch))
- Add an extra sanity check to improve error message [\#155](https://github.com/bazelbuild/bazel-watcher/pull/155) ([achew22](https://github.com/achew22))
- Upgrade to latest Bazel in Travis [\#154](https://github.com/bazelbuild/bazel-watcher/pull/154) ([achew22](https://github.com/achew22))
- Switch to python for workplace\_status script [\#153](https://github.com/bazelbuild/bazel-watcher/pull/153) ([achew22](https://github.com/achew22))
- Update dependency io\_bazel\_rules\_go to v0.16.1 [\#152](https://github.com/bazelbuild/bazel-watcher/pull/152) ([renovate[bot]](https://github.com/apps/renovate))
- Update dependency bazel\_gazelle to v0.15.0 [\#151](https://github.com/bazelbuild/bazel-watcher/pull/151) ([renovate[bot]](https://github.com/apps/renovate))
- Configure Renovate [\#150](https://github.com/bazelbuild/bazel-watcher/pull/150) ([renovate[bot]](https://github.com/apps/renovate))
- Fix wrong shell command [\#148](https://github.com/bazelbuild/bazel-watcher/pull/148) ([Xadeck](https://github.com/Xadeck))
- Support Bazel --keep\_going flag [\#147](https://github.com/bazelbuild/bazel-watcher/pull/147) ([schroederc](https://github.com/schroederc))

## [v0.6.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.6.0) (2018-10-08)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.5.0...v0.6.0)

**Closed issues:**

- Usage instructions are unclear [\#123](https://github.com/bazelbuild/bazel-watcher/issues/123)
- killing a server: tell the user to press ctrl-c a second time [\#122](https://github.com/bazelbuild/bazel-watcher/issues/122)
- Add explanation of where the ibazel binary lives to the README [\#114](https://github.com/bazelbuild/bazel-watcher/issues/114)

**Merged pull requests:**

- Support Bazel --output\_groups flag [\#145](https://github.com/bazelbuild/bazel-watcher/pull/145) ([schroederc](https://github.com/schroederc))
- Prompt for second SIGINT [\#142](https://github.com/bazelbuild/bazel-watcher/pull/142) ([dougkoch](https://github.com/dougkoch))
- Clean up README [\#141](https://github.com/bazelbuild/bazel-watcher/pull/141) ([dougkoch](https://github.com/dougkoch))
- Update bazelbuild/bazel-integration-testing to latest version [\#140](https://github.com/bazelbuild/bazel-watcher/pull/140) ([clintharrison](https://github.com/clintharrison))
- Add test reproducing vim backupcopy=no and kqueue failure [\#139](https://github.com/bazelbuild/bazel-watcher/pull/139) ([clintharrison](https://github.com/clintharrison))
- Update rules\_go and gazelle [\#137](https://github.com/bazelbuild/bazel-watcher/pull/137) ([kalbasit](https://github.com/kalbasit))

## [v0.5.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.5.0) (2018-08-15)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.4.0...v0.5.0)

**Closed issues:**

- undefined: js when trying to build it [\#132](https://github.com/bazelbuild/bazel-watcher/issues/132)

**Merged pull requests:**

- Ignore error about /tools/defaults/BUILD not existing [\#129](https://github.com/bazelbuild/bazel-watcher/pull/129) ([alexeagle](https://github.com/alexeagle))

## [v0.4.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.4.0) (2018-06-01)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.3.1...v0.4.0)

**Closed issues:**

- ERROR: Unrecognized option: --define=key=value [\#120](https://github.com/bazelbuild/bazel-watcher/issues/120)
- Your platform/architecture combination NaN is not yet supported. Windows 10 [\#116](https://github.com/bazelbuild/bazel-watcher/issues/116)
- Home for ibazel-benchmark-runner code? [\#111](https://github.com/bazelbuild/bazel-watcher/issues/111)
- Publish 0.3.0 to npm [\#110](https://github.com/bazelbuild/bazel-watcher/issues/110)
- Resource limit hit on OSX [\#101](https://github.com/bazelbuild/bazel-watcher/issues/101)
- Automatically apply buildozer commands from warnings [\#18](https://github.com/bazelbuild/bazel-watcher/issues/18)

**Merged pull requests:**

- Output Runner [\#134](https://github.com/bazelbuild/bazel-watcher/pull/134) ([borkaehw](https://github.com/borkaehw))
- No longer print current state name [\#128](https://github.com/bazelbuild/bazel-watcher/pull/128) ([mrmeku](https://github.com/mrmeku))
- automatically increase process ulimit [\#125](https://github.com/bazelbuild/bazel-watcher/pull/125) ([brendanjryan](https://github.com/brendanjryan))
- Bump Travis config to Bazel 0.10.0 [\#118](https://github.com/bazelbuild/bazel-watcher/pull/118) ([dougkoch](https://github.com/dougkoch))
- Optional config flag to be passed to bazel \(for CI use with --config=ci\) [\#113](https://github.com/bazelbuild/bazel-watcher/pull/113) ([gregmagolan](https://github.com/gregmagolan))

## [v0.3.1](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.3.1) (2018-01-19)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.3.0...v0.3.1)

**Merged pull requests:**

- Create a release script [\#115](https://github.com/bazelbuild/bazel-watcher/pull/115) ([achew22](https://github.com/achew22))

## [v0.3.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.3.0) (2017-12-22)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.2.0...v0.3.0)

**Closed issues:**

- process is left running after ctrl-c [\#99](https://github.com/bazelbuild/bazel-watcher/issues/99)
- e2e non-determinism after addition of lifecycle hooks [\#95](https://github.com/bazelbuild/bazel-watcher/issues/95)
- Error asm while building ibazel [\#91](https://github.com/bazelbuild/bazel-watcher/issues/91)
- broken when running from WORKSPACE/subdirectory [\#49](https://github.com/bazelbuild/bazel-watcher/issues/49)
- Make it possible to not restart the server and only rebuild the data dependencies of a job \(while leaving it running\) [\#29](https://github.com/bazelbuild/bazel-watcher/issues/29)
- ibazel occasionally stops running on changes [\#23](https://github.com/bazelbuild/bazel-watcher/issues/23)
- Leaking watchers will lead to resource exhaustion [\#11](https://github.com/bazelbuild/bazel-watcher/issues/11)

**Merged pull requests:**

- Move e2e output assertions into helper [\#103](https://github.com/bazelbuild/bazel-watcher/pull/103) ([achew22](https://github.com/achew22))
- Fixes shutdown of npm based iBazel [\#102](https://github.com/bazelbuild/bazel-watcher/pull/102) ([gregmagolan](https://github.com/gregmagolan))
- Profiler as lifecycle hook [\#98](https://github.com/bazelbuild/bazel-watcher/pull/98) ([gregmagolan](https://github.com/gregmagolan))
- Add the concept of lifecycle hooks [\#92](https://github.com/bazelbuild/bazel-watcher/pull/92) ([achew22](https://github.com/achew22))
- Create a directory per logical group of e2e test [\#89](https://github.com/bazelbuild/bazel-watcher/pull/89) ([achew22](https://github.com/achew22))
- Warn when using globally installed npm package [\#86](https://github.com/bazelbuild/bazel-watcher/pull/86) ([alexeagle](https://github.com/alexeagle))
- Upgrade bazel-integration-testing [\#85](https://github.com/bazelbuild/bazel-watcher/pull/85) ([achew22](https://github.com/achew22))
- Watch absolute paths [\#83](https://github.com/bazelbuild/bazel-watcher/pull/83) ([dougkoch](https://github.com/dougkoch))

## [v0.2.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.2.0) (2017-12-03)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.1.1...v0.2.0)

**Closed issues:**

- Versioning needed [\#77](https://github.com/bazelbuild/bazel-watcher/issues/77)
- avoid version problems for node users [\#45](https://github.com/bazelbuild/bazel-watcher/issues/45)

**Merged pull requests:**

- Upgrade rules\_go [\#84](https://github.com/bazelbuild/bazel-watcher/pull/84) ([achew22](https://github.com/achew22))
- Move e2e tests to bazel-integartion-testing. [\#82](https://github.com/bazelbuild/bazel-watcher/pull/82) ([achew22](https://github.com/achew22))
- Search for a local npm installation [\#74](https://github.com/bazelbuild/bazel-watcher/pull/74) ([dougkoch](https://github.com/dougkoch))
- Unwatch files no longer returned by query [\#71](https://github.com/bazelbuild/bazel-watcher/pull/71) ([dougkoch](https://github.com/dougkoch))

## [v0.1.1](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.1.1) (2017-12-01)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.0.1...v0.1.1)

**Closed issues:**

- bazel test fails locally and on TravisCI [\#72](https://github.com/bazelbuild/bazel-watcher/issues/72)
- PKGBUILD for Arch Linux [\#66](https://github.com/bazelbuild/bazel-watcher/issues/66)
- Crash when passing --profile [\#58](https://github.com/bazelbuild/bazel-watcher/issues/58)
- Communicate to a long-running process that it will get stdin notifications [\#57](https://github.com/bazelbuild/bazel-watcher/issues/57)
- ibazel writes IBAZEL\_BUILD\_COMPLETED FAILURE to stdin even for successful builds [\#56](https://github.com/bazelbuild/bazel-watcher/issues/56)
- Replace the "MAGIC" tag with something permanent; document it [\#55](https://github.com/bazelbuild/bazel-watcher/issues/55)
- query fails quietly, easy to overlook and hard to debug [\#54](https://github.com/bazelbuild/bazel-watcher/issues/54)
- Fix CI status badges [\#53](https://github.com/bazelbuild/bazel-watcher/issues/53)
- Sad error reporting: panic [\#43](https://github.com/bazelbuild/bazel-watcher/issues/43)
- Notify long-running processes when a build finishes [\#36](https://github.com/bazelbuild/bazel-watcher/issues/36)
- Request: add CLI option to set debounce delay value [\#31](https://github.com/bazelbuild/bazel-watcher/issues/31)
- When initial build fails, ibazel should exit [\#2](https://github.com/bazelbuild/bazel-watcher/issues/2)

**Merged pull requests:**

- Generate package.json and cross-build in Bazel [\#81](https://github.com/bazelbuild/bazel-watcher/pull/81) ([achew22](https://github.com/achew22))
- Ibazel version [\#79](https://github.com/bazelbuild/bazel-watcher/pull/79) ([gregmagolan](https://github.com/gregmagolan))
- Live reload [\#78](https://github.com/bazelbuild/bazel-watcher/pull/78) ([gregmagolan](https://github.com/gregmagolan))
- Bump Travis config to Bazel 0.8.0 [\#76](https://github.com/bazelbuild/bazel-watcher/pull/76) ([dougkoch](https://github.com/dougkoch))
- Fix and document notify tag [\#75](https://github.com/bazelbuild/bazel-watcher/pull/75) ([dougkoch](https://github.com/dougkoch))
- Update contributors [\#73](https://github.com/bazelbuild/bazel-watcher/pull/73) ([dougkoch](https://github.com/dougkoch))
- Notify of changes command: Correctly identify when a build completes [\#70](https://github.com/bazelbuild/bazel-watcher/pull/70) ([mrmeku](https://github.com/mrmeku))
- Specify that the exit codes are not an API [\#68](https://github.com/bazelbuild/bazel-watcher/pull/68) ([dougkoch](https://github.com/dougkoch))
- Update build status to point to new Jenkins path [\#64](https://github.com/bazelbuild/bazel-watcher/pull/64) ([achew22](https://github.com/achew22))
- Exit when query fails and pass through its stderr [\#63](https://github.com/bazelbuild/bazel-watcher/pull/63) ([dougkoch](https://github.com/dougkoch))
- Add flag for debounce delay [\#62](https://github.com/bazelbuild/bazel-watcher/pull/62) ([dougkoch](https://github.com/dougkoch))
- Update rules\_go to 0.7.0 [\#61](https://github.com/bazelbuild/bazel-watcher/pull/61) ([dougkoch](https://github.com/dougkoch))
- Notify command now watches on the correct tag. [\#60](https://github.com/bazelbuild/bazel-watcher/pull/60) ([achew22](https://github.com/achew22))
- Avoid triggering builds on file accesses. [\#59](https://github.com/bazelbuild/bazel-watcher/pull/59) ([mprobst](https://github.com/mprobst))
- Catch SIGTERM instead of SIGKILL [\#51](https://github.com/bazelbuild/bazel-watcher/pull/51) ([dougkoch](https://github.com/dougkoch))

## [v0.0.1](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.0.1) (2017-10-08)
**Closed issues:**

- Don't write all the watched file locations [\#40](https://github.com/bazelbuild/bazel-watcher/issues/40)
- ibazel run not forwarding ctrl-c to processes [\#32](https://github.com/bazelbuild/bazel-watcher/issues/32)
- Link error  "-pie and -r are incompatible" [\#26](https://github.com/bazelbuild/bazel-watcher/issues/26)
- Watch for new files in all packages [\#22](https://github.com/bazelbuild/bazel-watcher/issues/22)
- ibazel does not pass arguments to Bazel [\#13](https://github.com/bazelbuild/bazel-watcher/issues/13)
- ibazel doesn't support multiple target patterns even though it is documented to do so [\#12](https://github.com/bazelbuild/bazel-watcher/issues/12)
- Track files who are replaced \(by mv operations\) when being saved [\#9](https://github.com/bazelbuild/bazel-watcher/issues/9)
- ibazel doesn't watch changes in BUILD files that affect the build [\#8](https://github.com/bazelbuild/bazel-watcher/issues/8)
- Document how this compares with the `--watchfs` bazel option [\#5](https://github.com/bazelbuild/bazel-watcher/issues/5)
- file watcher triggers too much [\#4](https://github.com/bazelbuild/bazel-watcher/issues/4)
- ibazel doesn't watch changes in .bzl files that affect the build [\#3](https://github.com/bazelbuild/bazel-watcher/issues/3)

**Merged pull requests:**

- Add an end to end test [\#48](https://github.com/bazelbuild/bazel-watcher/pull/48) ([achew22](https://github.com/achew22))
- fix broken notification mode [\#47](https://github.com/bazelbuild/bazel-watcher/pull/47) ([alexeagle](https://github.com/alexeagle))
- Clean up the logging for successful watched files [\#46](https://github.com/bazelbuild/bazel-watcher/pull/46) ([mrmeku](https://github.com/mrmeku))
- Create script for generating npm releases [\#44](https://github.com/bazelbuild/bazel-watcher/pull/44) ([achew22](https://github.com/achew22))
- Include compiled .pb.go in the repo. [\#42](https://github.com/bazelbuild/bazel-watcher/pull/42) ([achew22](https://github.com/achew22))
- Detect the tag iblaze\_notify\_changes and use NotifyCommand instead [\#41](https://github.com/bazelbuild/bazel-watcher/pull/41) ([achew22](https://github.com/achew22))
- Create NotifyCommand mode [\#39](https://github.com/bazelbuild/bazel-watcher/pull/39) ([achew22](https://github.com/achew22))
- Restructure running commands [\#38](https://github.com/bazelbuild/bazel-watcher/pull/38) ([achew22](https://github.com/achew22))
- Added a command abstraction. [\#37](https://github.com/bazelbuild/bazel-watcher/pull/37) ([mrmeku](https://github.com/mrmeku))
- Listen for SIGINT/KILL and quit. [\#35](https://github.com/bazelbuild/bazel-watcher/pull/35) ([achew22](https://github.com/achew22))
- Use bazel run --script\_path [\#34](https://github.com/bazelbuild/bazel-watcher/pull/34) ([achew22](https://github.com/achew22))
- wait for process to end when killing it [\#33](https://github.com/bazelbuild/bazel-watcher/pull/33) ([jmhodges](https://github.com/jmhodges))
- Add a basic argument processor. [\#30](https://github.com/bazelbuild/bazel-watcher/pull/30) ([achew22](https://github.com/achew22))
- Kill preexisting when ibazel run detects changes [\#27](https://github.com/bazelbuild/bazel-watcher/pull/27) ([achew22](https://github.com/achew22))
- Update rules\_go to 0.5.4 [\#25](https://github.com/bazelbuild/bazel-watcher/pull/25) ([vladmos](https://github.com/vladmos))
- Upgrade to rules\_go 0.5.2 [\#19](https://github.com/bazelbuild/bazel-watcher/pull/19) ([mattmoor](https://github.com/mattmoor))
- Add travis-ci.com config [\#17](https://github.com/bazelbuild/bazel-watcher/pull/17) ([achew22](https://github.com/achew22))
- Refactor main.go into testable chunks [\#16](https://github.com/bazelbuild/bazel-watcher/pull/16) ([achew22](https://github.com/achew22))
- Support multiple target patterns [\#14](https://github.com/bazelbuild/bazel-watcher/pull/14) ([dougkoch](https://github.com/dougkoch))
- Track files replaced by move [\#10](https://github.com/bazelbuild/bazel-watcher/pull/10) ([dougkoch](https://github.com/dougkoch))
- Add a debounce of 100ms to ibazel actions. [\#7](https://github.com/bazelbuild/bazel-watcher/pull/7) ([achew22](https://github.com/achew22))
- Add documentation about --watchfs [\#6](https://github.com/bazelbuild/bazel-watcher/pull/6) ([achew22](https://github.com/achew22))
- Update Bazel URL [\#1](https://github.com/bazelbuild/bazel-watcher/pull/1) ([steren](https://github.com/steren))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*