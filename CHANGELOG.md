# Change Log

## [v0.13.1](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.13.1) (2020-05-18)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.13.0...v0.13.1)

**Implemented enhancements:**

- ibazel don't watch for changes in external repository [\#274](https://github.com/bazelbuild/bazel-watcher/issues/274)

**Fixed bugs:**

- Error reading config: open .bazel\_fix\_commands.json: no such file or directory [\#369](https://github.com/bazelbuild/bazel-watcher/issues/369)
- If the inner program crashes or dies, ibazel should restart it [\#300](https://github.com/bazelbuild/bazel-watcher/issues/300)

**Closed issues:**

- process not restarted on failure [\#385](https://github.com/bazelbuild/bazel-watcher/issues/385)
- ibazel with --override\_repository can't handle changes in the build tree [\#383](https://github.com/bazelbuild/bazel-watcher/issues/383)
- How do you enable JS serving? [\#358](https://github.com/bazelbuild/bazel-watcher/issues/358)
- Support multiple tartget [\#351](https://github.com/bazelbuild/bazel-watcher/issues/351)
- ibazel can't locate bazel if not in $PATH [\#341](https://github.com/bazelbuild/bazel-watcher/issues/341)

**Merged pull requests:**

- Update module golang/protobuf to v1.4.2 [\#389](https://github.com/bazelbuild/bazel-watcher/pull/389) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency bazel\_gazelle to v0.21.0 [\#387](https://github.com/bazelbuild/bazel-watcher/pull/387) ([renovate-bot](https://github.com/renovate-bot))
- Pass override repository flags to query commands and add e2e test [\#384](https://github.com/bazelbuild/bazel-watcher/pull/384) ([lewish](https://github.com/lewish))
- Allow `--build\_tag\_filters` and `--build\_tests\_only` to be specified  [\#382](https://github.com/bazelbuild/bazel-watcher/pull/382) ([devversion](https://github.com/devversion))
- Update module golang/protobuf to v1.4.1 [\#381](https://github.com/bazelbuild/bazel-watcher/pull/381) ([renovate-bot](https://github.com/renovate-bot))
- Corrections and cosmetic adjustments to README.md [\#380](https://github.com/bazelbuild/bazel-watcher/pull/380) ([trironkk](https://github.com/trironkk))
- Watch local repositories [\#379](https://github.com/bazelbuild/bazel-watcher/pull/379) ([lewish](https://github.com/lewish))
- Add an example client [\#378](https://github.com/bazelbuild/bazel-watcher/pull/378) ([achew22](https://github.com/achew22))
- Map rules\_go in gazelle directive [\#377](https://github.com/bazelbuild/bazel-watcher/pull/377) ([achew22](https://github.com/achew22))
- Restart subprocess on changes if it crashed [\#366](https://github.com/bazelbuild/bazel-watcher/pull/366) ([mrmeku](https://github.com/mrmeku))
- Update module bazelbuild/rules\_go to v0.23.1 [\#365](https://github.com/bazelbuild/bazel-watcher/pull/365) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency io\_bazel\_rules\_go to v0.22.4 [\#364](https://github.com/bazelbuild/bazel-watcher/pull/364) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency com\_github\_bazelbuild\_rules\_go to v0.22.4 [\#363](https://github.com/bazelbuild/bazel-watcher/pull/363) ([renovate-bot](https://github.com/renovate-bot))

## [v0.13.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.13.0) (2020-04-24)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.12.4...v0.13.0)

**Fixed bugs:**

- Reads .bazel\_fix\_commands.json from current directory [\#373](https://github.com/bazelbuild/bazel-watcher/issues/373)

**Closed issues:**

- Proposal:  use bazel's cquery in place of query for queryForSourceFiles checks [\#305](https://github.com/bazelbuild/bazel-watcher/issues/305)
- ibazel run on container\_image does not work [\#245](https://github.com/bazelbuild/bazel-watcher/issues/245)
- ibazel run crash doesn't shut down ts\_devserver [\#197](https://github.com/bazelbuild/bazel-watcher/issues/197)

**Merged pull requests:**

- Generating CHANGELOG.md for release v0.13.0 [\#376](https://github.com/bazelbuild/bazel-watcher/pull/376) ([achew22](https://github.com/achew22))
- Fix output\_runner to read from %WORKSPACE [\#375](https://github.com/bazelbuild/bazel-watcher/pull/375) ([achew22](https://github.com/achew22))
- Switch to cquery instead of query [\#374](https://github.com/bazelbuild/bazel-watcher/pull/374) ([achew22](https://github.com/achew22))
- Update module golang/protobuf to v1.4.0 [\#372](https://github.com/bazelbuild/bazel-watcher/pull/372) ([renovate-bot](https://github.com/renovate-bot))

## [v0.12.4](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.12.4) (2020-04-09)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.12.3...v0.12.4)

**Fixed bugs:**

- Full log not being shown on initial querying [\#217](https://github.com/bazelbuild/bazel-watcher/issues/217)
- Support for passing flags to bazel is limited and undocumented [\#126](https://github.com/bazelbuild/bazel-watcher/issues/126)

**Closed issues:**

- bazeliskNpmPath: should check if @bazel/bazelisk binary exists? [\#370](https://github.com/bazelbuild/bazel-watcher/issues/370)
- Bazelisk regression [\#352](https://github.com/bazelbuild/bazel-watcher/issues/352)
- \[Windows\] - Querying for files to watch... \)\) was unexpected at this time. Bazel query failed: exit status 255 [\#344](https://github.com/bazelbuild/bazel-watcher/issues/344)

**Merged pull requests:**

- Fix a lingering TODO [\#362](https://github.com/bazelbuild/bazel-watcher/pull/362) ([achew22](https://github.com/achew22))
- Stamp can be passed without a value [\#361](https://github.com/bazelbuild/bazel-watcher/pull/361) ([achew22](https://github.com/achew22))
- Update module golang/protobuf to v1.3.5 [\#359](https://github.com/bazelbuild/bazel-watcher/pull/359) ([renovate-bot](https://github.com/renovate-bot))
- Update module fsnotify/fsnotify to v1.4.9 [\#357](https://github.com/bazelbuild/bazel-watcher/pull/357) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency com\_github\_bazelbuild\_rules\_go to v0.22.1 [\#354](https://github.com/bazelbuild/bazel-watcher/pull/354) ([renovate-bot](https://github.com/renovate-bot))

## [v0.12.3](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.12.3) (2020-03-14)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.12.2...v0.12.3)

**Closed issues:**

- ibazel not connecting to bazel server on VPN [\#356](https://github.com/bazelbuild/bazel-watcher/issues/356)
- No artifacts uploaded for 0.12.0 and 0.11.2 [\#350](https://github.com/bazelbuild/bazel-watcher/issues/350)
- $TEST\_TMPDIR is deleted across reloads [\#323](https://github.com/bazelbuild/bazel-watcher/issues/323)
- Flag --incompatible\_no\_implicit\_file\_export will break Bazel watcher in a future Bazel release [\#319](https://github.com/bazelbuild/bazel-watcher/issues/319)
- Flag --incompatible\_no\_implicit\_file\_export will break Bazel watcher in Bazel 1.2.1 [\#316](https://github.com/bazelbuild/bazel-watcher/issues/316)
- Bash trap works with bazel but not with ibazel [\#291](https://github.com/bazelbuild/bazel-watcher/issues/291)

**Merged pull requests:**

- Fix stamping regression caused by latest rules\_go [\#360](https://github.com/bazelbuild/bazel-watcher/pull/360) ([achew22](https://github.com/achew22))
- Update module golang/protobuf to v1.3.4 [\#355](https://github.com/bazelbuild/bazel-watcher/pull/355) ([renovate-bot](https://github.com/renovate-bot))
- Enable run\_output by default [\#349](https://github.com/bazelbuild/bazel-watcher/pull/349) ([achew22](https://github.com/achew22))
- Update module bazelbuild/rules\_go to v0.22.1 [\#336](https://github.com/bazelbuild/bazel-watcher/pull/336) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency io\_bazel\_rules\_go to v0.22.1 [\#335](https://github.com/bazelbuild/bazel-watcher/pull/335) ([renovate-bot](https://github.com/renovate-bot))

## [v0.12.2](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.12.2) (2020-02-24)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.12.1...v0.12.2)

## [v0.12.1](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.12.1) (2020-02-24)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.12.0...v0.12.1)

## [v0.12.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.12.0) (2020-02-19)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.11.2...v0.12.0)

**Closed issues:**

- bazel watcher does not work with `run -c` [\#347](https://github.com/bazelbuild/bazel-watcher/issues/347)
- Try resolve bazel binary from @bazel/bazelisk installed locally [\#339](https://github.com/bazelbuild/bazel-watcher/issues/339)

**Merged pull requests:**

- Support bazelisk [\#346](https://github.com/bazelbuild/bazel-watcher/pull/346) ([zoidbergwill](https://github.com/zoidbergwill))

## [v0.11.2](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.11.2) (2020-02-14)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.11.1...v0.11.2)

**Closed issues:**

- Querying for files to watch... \)\) was unexpected at this time. Bazel query failed: exit status 255 [\#342](https://github.com/bazelbuild/bazel-watcher/issues/342)
- Flag --incompatible\_load\_proto\_rules\_from\_bzl will break Bazel watcher in Bazel 1.2.1 [\#317](https://github.com/bazelbuild/bazel-watcher/issues/317)

**Merged pull requests:**

- Adds passthrough for bazel flags "--compilation\_mode" and "-c" [\#348](https://github.com/bazelbuild/bazel-watcher/pull/348) ([Jdban](https://github.com/Jdban))
- Update dependency com\_github\_bazelbuild\_rules\_go to v0.21.3 [\#345](https://github.com/bazelbuild/bazel-watcher/pull/345) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency bazel\_gazelle to v0.20.0 [\#343](https://github.com/bazelbuild/bazel-watcher/pull/343) ([renovate-bot](https://github.com/renovate-bot))
- support WORKSPACE.bazel file [\#340](https://github.com/bazelbuild/bazel-watcher/pull/340) ([alexeagle](https://github.com/alexeagle))
- Update rules\_proto commit hash to f6b8d89 [\#338](https://github.com/bazelbuild/bazel-watcher/pull/338) ([renovate-bot](https://github.com/renovate-bot))
- Update module golang/protobuf to v1.3.3 [\#337](https://github.com/bazelbuild/bazel-watcher/pull/337) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency com\_github\_bazelbuild\_rules\_go to v0.21.2 [\#334](https://github.com/bazelbuild/bazel-watcher/pull/334) ([renovate-bot](https://github.com/renovate-bot))
- Update module bazelbuild/rules\_go to v0.21.0 [\#332](https://github.com/bazelbuild/bazel-watcher/pull/332) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency com\_github\_bazelbuild\_rules\_go to v0.21.0 [\#331](https://github.com/bazelbuild/bazel-watcher/pull/331) ([renovate-bot](https://github.com/renovate-bot))
- Update module bazelbuild/rules\_go to v0.20.4 [\#328](https://github.com/bazelbuild/bazel-watcher/pull/328) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency io\_bazel\_rules\_go to v0.21.0 [\#327](https://github.com/bazelbuild/bazel-watcher/pull/327) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency com\_github\_bazelbuild\_rules\_go to v0.20.4 [\#326](https://github.com/bazelbuild/bazel-watcher/pull/326) ([renovate-bot](https://github.com/renovate-bot))
- Update rules\_proto commit hash to d7666ec [\#325](https://github.com/bazelbuild/bazel-watcher/pull/325) ([renovate-bot](https://github.com/renovate-bot))
- Load proto from @rules\_proto [\#321](https://github.com/bazelbuild/bazel-watcher/pull/321) ([achew22](https://github.com/achew22))

## [v0.11.1](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.11.1) (2020-01-07)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.11.0...v0.11.1)

**Closed issues:**

- parallel run support [\#320](https://github.com/bazelbuild/bazel-watcher/issues/320)
- Give output runner an option to exit early [\#260](https://github.com/bazelbuild/bazel-watcher/issues/260)
- Make livereload usable for custom rules  [\#248](https://github.com/bazelbuild/bazel-watcher/issues/248)
- Add option to ignore files to be watched [\#244](https://github.com/bazelbuild/bazel-watcher/issues/244)
- File changes beyond the first are not detected on Windows 10 [\#236](https://github.com/bazelbuild/bazel-watcher/issues/236)
- ibazel\_notify\_changes not working for tests [\#184](https://github.com/bazelbuild/bazel-watcher/issues/184)
- Watching for files being added [\#135](https://github.com/bazelbuild/bazel-watcher/issues/135)
- -log\_to\_file redirects bazel's stderr in addition to ibazel's [\#124](https://github.com/bazelbuild/bazel-watcher/issues/124)

**Merged pull requests:**

- Add --nocache\_test\_results to overrideableBazelFlags [\#324](https://github.com/bazelbuild/bazel-watcher/pull/324) ([mariusgrigoriu](https://github.com/mariusgrigoriu))
- Remove tap-bin from Home-brew formula. [\#318](https://github.com/bazelbuild/bazel-watcher/pull/318) ([BooneJS](https://github.com/BooneJS))

## [v0.11.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.11.0) (2019-12-17)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.10.3...v0.11.0)

**Closed issues:**

- runing != running [\#313](https://github.com/bazelbuild/bazel-watcher/issues/313)
- v0.10.3 fails to build due to error in bazel-integration-testing//tools:common.bzl dependency [\#297](https://github.com/bazelbuild/bazel-watcher/issues/297)
- Clarification on tag/release workflow [\#296](https://github.com/bazelbuild/bazel-watcher/issues/296)
- Color output in console can conflict with output\_runner regex [\#263](https://github.com/bazelbuild/bazel-watcher/issues/263)

**Merged pull requests:**

- Update dependencies [\#315](https://github.com/bazelbuild/bazel-watcher/pull/315) ([achew22](https://github.com/achew22))
- Verbize run to running not runing [\#314](https://github.com/bazelbuild/bazel-watcher/pull/314) ([achew22](https://github.com/achew22))
- Fix NPM release script [\#312](https://github.com/bazelbuild/bazel-watcher/pull/312) ([achew22](https://github.com/achew22))
- Update dependency com\_github\_bazelbuild\_rules\_go to v0.20.3 [\#309](https://github.com/bazelbuild/bazel-watcher/pull/309) ([renovate-bot](https://github.com/renovate-bot))
- Update rules\_proto commit hash to 2c04683 [\#308](https://github.com/bazelbuild/bazel-watcher/pull/308) ([renovate-bot](https://github.com/renovate-bot))
- Update rules\_proto commit hash to f6c112f [\#307](https://github.com/bazelbuild/bazel-watcher/pull/307) ([renovate-bot](https://github.com/renovate-bot))
- Update to latest rules\_go [\#303](https://github.com/bazelbuild/bazel-watcher/pull/303) ([achew22](https://github.com/achew22))
- Add a logging module [\#302](https://github.com/bazelbuild/bazel-watcher/pull/302) ([achew22](https://github.com/achew22))
- Remove ANSI codes from output before matching [\#299](https://github.com/bazelbuild/bazel-watcher/pull/299) ([DavidANeil](https://github.com/DavidANeil))
- Add installation instructions for Arch Linux [\#298](https://github.com/bazelbuild/bazel-watcher/pull/298) ([sudoforge](https://github.com/sudoforge))
- Update com\_github\_bazelbuild\_bazel\_integration\_testing commit hash to cd4f16d [\#268](https://github.com/bazelbuild/bazel-watcher/pull/268) ([renovate-bot](https://github.com/renovate-bot))

## [v0.10.3](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.10.3) (2019-11-01)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.10.2...v0.10.3)

**Closed issues:**

- error setting higher file descriptor limit for this process: invalid argument [\#285](https://github.com/bazelbuild/bazel-watcher/issues/285)
- Add support for bazelrc flag [\#278](https://github.com/bazelbuild/bazel-watcher/issues/278)
- run command should also forward recognized flags [\#269](https://github.com/bazelbuild/bazel-watcher/issues/269)
- --test\_env argument not in allowed list [\#256](https://github.com/bazelbuild/bazel-watcher/issues/256)
- ibazel can't locate bazel if not in $PATH [\#252](https://github.com/bazelbuild/bazel-watcher/issues/252)
- Add Windows support [\#105](https://github.com/bazelbuild/bazel-watcher/issues/105)

**Merged pull requests:**

- add --copt= to overrideable flags [\#295](https://github.com/bazelbuild/bazel-watcher/pull/295) ([girtsf](https://github.com/girtsf))
- Update dependency com\_github\_bazelbuild\_rules\_go to v0.20.2 [\#294](https://github.com/bazelbuild/bazel-watcher/pull/294) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency bazel\_skylib to v1.0.2 [\#293](https://github.com/bazelbuild/bazel-watcher/pull/293) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency bazel\_skylib to v1.0.1 [\#292](https://github.com/bazelbuild/bazel-watcher/pull/292) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency bazel\_skylib to v1 [\#290](https://github.com/bazelbuild/bazel-watcher/pull/290) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency com\_github\_bazelbuild\_rules\_go to v0.19.5 [\#289](https://github.com/bazelbuild/bazel-watcher/pull/289) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency com\_github\_bazelbuild\_rules\_go to v0.19.4 [\#288](https://github.com/bazelbuild/bazel-watcher/pull/288) ([renovate-bot](https://github.com/renovate-bot))
- Most bazel flags take arguments with an = [\#287](https://github.com/bazelbuild/bazel-watcher/pull/287) ([achew22](https://github.com/achew22))
- Work around getrlimit syscall error on darwin [\#286](https://github.com/bazelbuild/bazel-watcher/pull/286) ([jeremyschlatter](https://github.com/jeremyschlatter))
- Update dependency com\_github\_gorilla\_websocket to v1.4.1 [\#284](https://github.com/bazelbuild/bazel-watcher/pull/284) ([renovate-bot](https://github.com/renovate-bot))
- Set IBAZEL=true in process env [\#282](https://github.com/bazelbuild/bazel-watcher/pull/282) ([statik](https://github.com/statik))
- Add --test\_tag\_filters= to overrideable flags [\#281](https://github.com/bazelbuild/bazel-watcher/pull/281) ([whilp](https://github.com/whilp))
- Add support for startup option bazelrc [\#280](https://github.com/bazelbuild/bazel-watcher/pull/280) ([libsamek](https://github.com/libsamek))
- Add documentation for the output runner [\#279](https://github.com/bazelbuild/bazel-watcher/pull/279) ([DavidANeil](https://github.com/DavidANeil))
- feat: locate Bazel binary in sibling package when installed via npm [\#275](https://github.com/bazelbuild/bazel-watcher/pull/275) ([alexeagle](https://github.com/alexeagle))
- Update dependency bazel\_skylib to v0.9.0 [\#273](https://github.com/bazelbuild/bazel-watcher/pull/273) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency com\_github\_golang\_protobuf to v1.3.2 [\#272](https://github.com/bazelbuild/bazel-watcher/pull/272) ([renovate-bot](https://github.com/renovate-bot))
- Pass bazel args for run command. [\#271](https://github.com/bazelbuild/bazel-watcher/pull/271) ([zoido](https://github.com/zoido))
- Update dependency io\_bazel\_rules\_go to v0.18.6 [\#267](https://github.com/bazelbuild/bazel-watcher/pull/267) ([renovate-bot](https://github.com/renovate-bot))
- Update com\_github\_bazelbuild\_bazel\_integration\_testing commit hash to 2ce0893 [\#266](https://github.com/bazelbuild/bazel-watcher/pull/266) ([renovate-bot](https://github.com/renovate-bot))
- Notify listening processes when a build starts. [\#265](https://github.com/bazelbuild/bazel-watcher/pull/265) ([DavidANeil](https://github.com/DavidANeil))
- Update com\_github\_bazelbuild\_bazel\_integration\_testing commit hash to b8f2053 [\#264](https://github.com/bazelbuild/bazel-watcher/pull/264) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency io\_bazel\_rules\_go to v0.18.5 [\#258](https://github.com/bazelbuild/bazel-watcher/pull/258) ([renovate-bot](https://github.com/renovate-bot))
- Add more Bazel flags [\#257](https://github.com/bazelbuild/bazel-watcher/pull/257) ([aaliddell](https://github.com/aaliddell))
- Update dependency io\_bazel\_rules\_go to v0.18.4 [\#255](https://github.com/bazelbuild/bazel-watcher/pull/255) ([renovate-bot](https://github.com/renovate-bot))
- Update com\_github\_bazelbuild\_bazel\_integration\_testing commit hash to 13a7d51 [\#254](https://github.com/bazelbuild/bazel-watcher/pull/254) ([renovate-bot](https://github.com/renovate-bot))

## [v0.10.2](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.10.2) (2019-05-01)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.10.1...v0.10.2)

**Merged pull requests:**

- Fix windows binary resolution on node [\#253](https://github.com/bazelbuild/bazel-watcher/pull/253) ([rerion](https://github.com/rerion))

## [v0.10.1](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.10.1) (2019-04-17)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.10.0...v0.10.1)

**Closed issues:**

- --define flag not passed on properly [\#242](https://github.com/bazelbuild/bazel-watcher/issues/242)

**Merged pull requests:**

- Update dependency io\_bazel\_rules\_go to v0.18.3 [\#250](https://github.com/bazelbuild/bazel-watcher/pull/250) ([renovate-bot](https://github.com/renovate-bot))
- refactor\(@bazel/ibazel\): export getNativeBinary [\#249](https://github.com/bazelbuild/bazel-watcher/pull/249) ([kyliau](https://github.com/kyliau))
- Update com\_github\_bazelbuild\_bazel\_integration\_testing commit hash to 83c0657 [\#247](https://github.com/bazelbuild/bazel-watcher/pull/247) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency io\_bazel\_rules\_go to v0.18.2 [\#246](https://github.com/bazelbuild/bazel-watcher/pull/246) ([renovate-bot](https://github.com/renovate-bot))
- Pass --define flag to Bazel [\#243](https://github.com/bazelbuild/bazel-watcher/pull/243) ([achew22](https://github.com/achew22))
- Update dependency bazel\_skylib to v0.8.0 [\#241](https://github.com/bazelbuild/bazel-watcher/pull/241) ([renovate-bot](https://github.com/renovate-bot))
- Update com\_github\_bazelbuild\_bazel\_integration\_testing commit hash to 5239780 [\#240](https://github.com/bazelbuild/bazel-watcher/pull/240) ([renovate-bot](https://github.com/renovate-bot))

## [v0.10.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.10.0) (2019-03-17)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.9.1...v0.10.0)

**Closed issues:**

- Bazel incompatible changes [\#224](https://github.com/bazelbuild/bazel-watcher/issues/224)
- Get all tests passing with --incompatible\_disable\_legacy\_cc\_provider for Bazel 0.25.0 [\#223](https://github.com/bazelbuild/bazel-watcher/issues/223)
- Fix the e2e tests in Windows [\#213](https://github.com/bazelbuild/bazel-watcher/issues/213)
- Bazel incompatible changes [\#207](https://github.com/bazelbuild/bazel-watcher/issues/207)
- ibazel 0.9.0 release assets [\#205](https://github.com/bazelbuild/bazel-watcher/issues/205)
- Add support for http\_archive [\#127](https://github.com/bazelbuild/bazel-watcher/issues/127)
- data race in ibazel [\#52](https://github.com/bazelbuild/bazel-watcher/issues/52)

**Merged pull requests:**

- Release Windows to Github [\#239](https://github.com/bazelbuild/bazel-watcher/pull/239) ([achew22](https://github.com/achew22))
- Release windows binaries in release scripts [\#238](https://github.com/bazelbuild/bazel-watcher/pull/238) ([achew22](https://github.com/achew22))
- Upgrade to checked in version of bazel integration [\#237](https://github.com/bazelbuild/bazel-watcher/pull/237) ([achew22](https://github.com/achew22))
- Update dependency io\_bazel\_rules\_go to v0.18.1 [\#234](https://github.com/bazelbuild/bazel-watcher/pull/234) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency com\_github\_golang\_protobuf to v1.3.1 [\#233](https://github.com/bazelbuild/bazel-watcher/pull/233) ([renovate-bot](https://github.com/renovate-bot))
- Enable e2e tests on Windows [\#232](https://github.com/bazelbuild/bazel-watcher/pull/232) ([meteorcloudy](https://github.com/meteorcloudy))
- Added newline to 'Starting...' text [\#231](https://github.com/bazelbuild/bazel-watcher/pull/231) ([zaucy](https://github.com/zaucy))
- Update com\_github\_bazelbuild\_bazel\_integration\_testing commit hash to c0a0423 [\#230](https://github.com/bazelbuild/bazel-watcher/pull/230) ([renovate-bot](https://github.com/renovate-bot))
- Update com\_github\_bazelbuild\_bazel\_integration\_testing commit hash to 5442699 [\#229](https://github.com/bazelbuild/bazel-watcher/pull/229) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency io\_bazel\_rules\_go to v0.18.0 [\#228](https://github.com/bazelbuild/bazel-watcher/pull/228) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency com\_github\_golang\_protobuf to v1.3.0 [\#227](https://github.com/bazelbuild/bazel-watcher/pull/227) ([renovate-bot](https://github.com/renovate-bot))
- Update dependency bazel\_gazelle to v0.17.0 [\#226](https://github.com/bazelbuild/bazel-watcher/pull/226) ([renovate-bot](https://github.com/renovate-bot))
- Update com\_github\_bazelbuild\_bazel\_integration\_testing commit hash to 9af28bd [\#225](https://github.com/bazelbuild/bazel-watcher/pull/225) ([renovate-bot](https://github.com/renovate-bot))
- Delete .travis.yml [\#222](https://github.com/bazelbuild/bazel-watcher/pull/222) ([achew22](https://github.com/achew22))
- Update com\_github\_bazelbuild\_bazel\_integration\_testing commit hash to 2cc1cd2 [\#221](https://github.com/bazelbuild/bazel-watcher/pull/221) ([renovate-bot](https://github.com/renovate-bot))
- Upgrade bazel-integration-testing [\#220](https://github.com/bazelbuild/bazel-watcher/pull/220) ([achew22](https://github.com/achew22))
- Update org\_golang\_x\_sys commit hash to cc5685c [\#219](https://github.com/bazelbuild/bazel-watcher/pull/219) ([renovate[bot]](https://github.com/apps/renovate))
- Update org\_golang\_x\_sys commit hash to ec7b60b [\#218](https://github.com/bazelbuild/bazel-watcher/pull/218) ([renovate[bot]](https://github.com/apps/renovate))

## [v0.9.1](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.9.1) (2019-02-17)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.9.0...v0.9.1)

**Closed issues:**

- iBazel should tell tell Bazel what files changed [\#201](https://github.com/bazelbuild/bazel-watcher/issues/201)

**Merged pull requests:**

- Add --test\_filter flag to ibazel [\#216](https://github.com/bazelbuild/bazel-watcher/pull/216) ([achew22](https://github.com/achew22))
- Upgrade to rules\_go 0.17.0 [\#215](https://github.com/bazelbuild/bazel-watcher/pull/215) ([achew22](https://github.com/achew22))
- Update dependency bazel\_skylib to v0.7.0 [\#214](https://github.com/bazelbuild/bazel-watcher/pull/214) ([renovate[bot]](https://github.com/apps/renovate))
- Create CODEOWNERS [\#212](https://github.com/bazelbuild/bazel-watcher/pull/212) ([dslomov](https://github.com/dslomov))
- Update org\_golang\_x\_sys commit hash to d0b11bd [\#210](https://github.com/bazelbuild/bazel-watcher/pull/210) ([renovate[bot]](https://github.com/apps/renovate))
- Update dependency bazel\_skylib to v0.6.0 [\#209](https://github.com/bazelbuild/bazel-watcher/pull/209) ([renovate[bot]](https://github.com/apps/renovate))
- Update dependency io\_bazel\_rules\_go to v0.16.6 [\#208](https://github.com/bazelbuild/bazel-watcher/pull/208) ([renovate[bot]](https://github.com/apps/renovate))
- Update bazel\_integration\_testing [\#204](https://github.com/bazelbuild/bazel-watcher/pull/204) ([meteorcloudy](https://github.com/meteorcloudy))
- Update dependency bazel\_gazelle to v0.16.0 [\#203](https://github.com/bazelbuild/bazel-watcher/pull/203) ([renovate[bot]](https://github.com/apps/renovate))
- Update dependency io\_bazel\_rules\_go to v0.16.5 [\#202](https://github.com/bazelbuild/bazel-watcher/pull/202) ([renovate[bot]](https://github.com/apps/renovate))
- Update org\_golang\_x\_sys commit hash to b907332 [\#199](https://github.com/bazelbuild/bazel-watcher/pull/199) ([renovate[bot]](https://github.com/apps/renovate))
- Add Windows support [\#144](https://github.com/bazelbuild/bazel-watcher/pull/144) ([jchv](https://github.com/jchv))

## [v0.9.0](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/tree/v0.9.0) (2018-12-07)
[Full Changelog](https://github.com/bazelbuild/bazel-watcher/bazelbuild/bazel-watcher/compare/v0.8.2...v0.9.0)

**Closed issues:**

- How do we test if output\_runner is reading the output correctly [\#196](https://github.com/bazelbuild/bazel-watcher/issues/196)
- ibazel should not crash if a build file change results in a syntax error [\#194](https://github.com/bazelbuild/bazel-watcher/issues/194)

**Merged pull requests:**

- Stay alive after source query failure [\#200](https://github.com/bazelbuild/bazel-watcher/pull/200) ([dougkoch](https://github.com/dougkoch))
- Add e2e test to verify output runner functionality [\#198](https://github.com/bazelbuild/bazel-watcher/pull/198) ([achew22](https://github.com/achew22))
- Fix the output runner [\#195](https://github.com/bazelbuild/bazel-watcher/pull/195) ([borkaehw](https://github.com/borkaehw))
- Update org\_golang\_x\_sys commit hash to a5c9d58 [\#193](https://github.com/bazelbuild/bazel-watcher/pull/193) ([renovate[bot]](https://github.com/apps/renovate))
- Add recent contributors [\#192](https://github.com/bazelbuild/bazel-watcher/pull/192) ([dougkoch](https://github.com/dougkoch))
- Run gofmt and fix typos [\#191](https://github.com/bazelbuild/bazel-watcher/pull/191) ([dougkoch](https://github.com/dougkoch))
- Update github publish script to build ghr [\#190](https://github.com/bazelbuild/bazel-watcher/pull/190) ([achew22](https://github.com/achew22))

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