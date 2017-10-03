#!/usr/bin/env node
// Copyright 2017 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
'use strict';

// This package inspired by
// https://github.com/angular/clang-format/blob/master/index.js
const os = require('os');
const path = require('path');
const spawn = require('child_process').spawn;

function main(args) {
  const nativeBinary =
      path.join(__dirname, 'bin', os.platform() + '_' + os.arch(), 'ibazel');
  if (os.platform() === 'darwin') {
    var nativeProcess = spawn(nativeBinary, args, {stdio: 'inherit'});
    nativeProcess.on('close', e => process.exitCode = e);
  } else {
    console.error(`FATAL: Platform ${os.platform()} not yet supported
    Follow install instructions at https://github.com/bazelbuild/bazel-watcher/blob/master/README.md`);
    return Promise.resolve(1);
  }
  process.exitCode = 0;
}

if (require.main === module) main(process.argv.slice(2));