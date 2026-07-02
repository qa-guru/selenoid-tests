import assert from "node:assert/strict";
import test from "node:test";
import { resolveDisplayName } from "./junit-to-allure.mjs";

test("resolveDisplayName uses t.Run subtest title after slash", () => {
  assert.equal(
    resolveDisplayName(
      "TestBadJsonFormat/POST /wd/hub/session rejects malformed JSON body",
      "github.com/aerokube/selenoid.TestBadJsonFormat/POST /wd/hub/session rejects malformed JSON body",
    ),
    "POST /wd/hub/session rejects malformed JSON body",
  );
});

test("resolveDisplayName uses fullName when legacy name has no subtest slash", () => {
  assert.equal(resolveDisplayName("TestBrowserNotFound", "github.com/aerokube/selenoid.TestBrowserNotFound"), "TestBrowserNotFound");
});

test("resolveDisplayName reads subtest from dotted fullName when name is absent", () => {
  assert.equal(
    resolveDisplayName(undefined, "github.com/aerokube/selenoid.TestBadJsonFormat/POST /wd/hub/session rejects malformed JSON body"),
    "POST /wd/hub/session rejects malformed JSON body",
  );
});

test("resolveDisplayName prefers name over fullName when both have slash", () => {
  assert.equal(resolveDisplayName("TestFoo/bar case", "pkg.TestFoo/bar case"), "bar case");
});
