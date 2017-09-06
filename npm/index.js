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
  const arch = {
    'x64': 'amd64',
  }[os.arch()];
  // Filter the platform based on the platforms that are build/included.
  const platform = {
    'darwin': 'darwin',
    'linux': 'linux',
    'windows': 'windows',
  }[os.platform()];

  if (arch == undefined || platform == undefined) {
    console.error(`FATAL: Your platform/architecture combination ${
        os.platform() - os.arch()} is not yet supported.
    Follow install instructions at https://github.com/bazelbuild/bazel-watcher/blob/master/README.md to compile for your system.`);
    return Promise.resolve(1);
  }

  const binary = path.join(__dirname, 'bin', `${platform}_${arch}`, 'ibazel');
  const ibazel = spawn(binary, args, {stdio: 'inherit'});
  ibazel.on('close', e => process.exitCode = e);
}

if (require.main === module) {
  main(process.argv.slice(2));
}
