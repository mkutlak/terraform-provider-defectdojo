## 0.1.0-rc1

BREAKING CHANGES:
 - Upgraded to Go 1.25 and updated all dependencies.
 - Replaced external `defect-dojo-client-go` with in-repo generated client (`internal/ddclient/`) via oapi-codegen v2, targeting DefectDojo API v2.54.3.
 - Upgraded Terraform Plugin Framework to latest version.

FEATURES:
 - **Infrastructure resources & data sources:**
   - `defectdojo_development_environment`
   - `defectdojo_regulation`
   - `defectdojo_tool_type`
   - `defectdojo_tool_configuration`
   - `defectdojo_sla_configuration`
   - `defectdojo_note_type`
   - `defectdojo_network_location`
   - `defectdojo_language_type`
 - **Security & Access Control resources & data sources:**
   - `defectdojo_user`
   - `defectdojo_user_contact_info`
   - `defectdojo_dojo_group`
   - `defectdojo_dojo_group_member`
   - `defectdojo_global_role`
   - `defectdojo_product_member`
   - `defectdojo_product_group`
   - `defectdojo_product_type_member`
   - `defectdojo_product_type_group`
   - `defectdojo_credential`
 - **Vulnerability Management resources & data sources:**
   - `defectdojo_engagement`
   - `defectdojo_engagement_preset`
   - `defectdojo_test`
   - `defectdojo_finding_template`
   - `defectdojo_endpoint`
 - **Integration resources & data sources:**
   - `defectdojo_jira_instance`
   - `defectdojo_tool_product_settings`
   - `defectdojo_product_api_scan_configuration`
   - `defectdojo_credential_mapping`
   - `defectdojo_risk_acceptance`
   - `defectdojo_notification_webhook`
   - `defectdojo_asset_group`
 - Add `defectdojo_jira_product_configuration` data source (previously resource-only).
 - Add new fields to `defectdojo_product`: `disable_sla_breach_notifications`, `enable_product_tag_inheritance`, `sla_configuration`.
 - Extended reflection engine to support `time.Time`, `Date`, `float64`, `int32`, non-pointer slices, and named string types.

IMPROVEMENTS:
 - Add Docker Compose setup for local DefectDojo acceptance testing.
 - Add `testAccPreCheck()` with environment variable validation.
 - Replace deprecated `resource.UniqueId()` in all test files.
 - Remove commented-out debug code from reflection engine.
 - Update CI/CD workflows: Go 1.25, actions v6, Terraform matrix 1.8/1.9/1.10.

## [1.1.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v1.0.1...v1.1.0) (2026-08-20)


### Features

* support clearing optional attributes removed from configuration ([e13db15](https://github.com/mkutlak/terraform-provider-defectdojo/commit/e13db1596c319fd1eb421b9a34bc3ee62b2ec639))


### Bug Fixes

* address review findings on tags, notifications and clearing ([25af402](https://github.com/mkutlak/terraform-provider-defectdojo/commit/25af402c778455983bf28e128b23ade4a4d1eede))
* enforce the tag grammar on engagement, test and url ([8afc1b2](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8afc1b255b205d395b33e89e2ee1c7b13461bc60))
* **engagement_preset:** require title, default scope to an empty string ([1cd0b4c](https://github.com/mkutlak/terraform-provider-defectdojo/commit/1cd0b4c82cce81d3c7abbee60ecb311670cf8887))
* **jira_instance:** stop overwriting the configured password on refresh ([34b72cf](https://github.com/mkutlak/terraform-provider-defectdojo/commit/34b72cffb87b9649cb0de48a235a0769b0bec80c))
* **location_product:** accept the empty relationship DefectDojo stores ([0e276c2](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0e276c2d5f3ecd13a94d0513c6200f57fc7cc64e))
* mark server-filled attributes as computed ([2a32792](https://github.com/mkutlak/terraform-provider-defectdojo/commit/2a32792f6e2e0ee8ba7e4c69e44339375095cdde))
* name the configuration attribute in the clearing diagnostics ([7aea1da](https://github.com/mkutlak/terraform-provider-defectdojo/commit/7aea1dac4f3373e1922c75d271d6ab4d1e23eb9d))
* **product:** accept single-character and multi-line descriptions ([8b47e2b](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8b47e2b78b72c3d35eb42f2bfb960e780234e67d))
* **product:** preserve the configured revenue literal ([6f6fdd1](https://github.com/mkutlak/terraform-provider-defectdojo/commit/6f6fdd1c451b7c2c6b6e8484f3b6265f12a24ee5))
* **product:** reject revenue literals that are not decimal numbers ([986e45b](https://github.com/mkutlak/terraform-provider-defectdojo/commit/986e45b96f803f169697367e5ef400be3fef37eb))
* send the clearing PATCH for attributes whose value was zero ([a132409](https://github.com/mkutlak/terraform-provider-defectdojo/commit/a1324096edae3520ae98af39809eb6f01cfd01e0))
* **tags:** reject case-colliding tags and drop duplicates from the read ([22eedcb](https://github.com/mkutlak/terraform-provider-defectdojo/commit/22eedcbcb83897d01c5c256cee61639f3cb2b6cc))
* **url:** preserve the configured host spelling ([8c0d75a](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8c0d75aec2bbddfece1222131e84a0bce0cc6729))


### Documentation

* record that product tag inheritance breaks Terraform-managed children ([cb4d666](https://github.com/mkutlak/terraform-provider-defectdojo/commit/cb4d6668212398ed4303906eda66cbdd855b4686))
* rewrite the schema descriptions in Simplified Technical English ([d5b4a33](https://github.com/mkutlak/terraform-provider-defectdojo/commit/d5b4a3351bff291afc73b9c656509aafc09379d2))
* **sla_configuration:** record that the day counts cannot be unset ([b365b74](https://github.com/mkutlak/terraform-provider-defectdojo/commit/b365b74fbd3e9dfd29d7d683be8b213a2b7d937c))
* **tags:** document what the tags attribute actually enforces ([7cb20f9](https://github.com/mkutlak/terraform-provider-defectdojo/commit/7cb20f92f9bb9b68d46bd70d36e62cbe8eec0362))


### CI/CD

* collect the OpenAPI spec before the acceptance runs ([cd55049](https://github.com/mkutlak/terraform-provider-defectdojo/commit/cd550494e96ef0f1f03c8ccaae4254b3ec5acf44))


### Code Refactoring

* assert reflect values without boxing them first ([7dd4f0f](https://github.com/mkutlak/terraform-provider-defectdojo/commit/7dd4f0ffca3db43243d0de11b8eb3df100b6dcfd))
* derive the clearing PATCH from the generated client ([d969118](https://github.com/mkutlak/terraform-provider-defectdojo/commit/d969118a3968c6684c3c45163008a51722975714))
* share one literal-preservation helper across the three formats ([3e4dc48](https://github.com/mkutlak/terraform-provider-defectdojo/commit/3e4dc48844a84459dc93dca707ffb84c2ecda950))
* surface dropped values as diagnostics instead of log warnings ([ba45b71](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ba45b71455f8b70a1ab56785205ad90149e380c0))
* **tags:** cut the commentary back to what the server actually does ([72bef70](https://github.com/mkutlak/terraform-provider-defectdojo/commit/72bef707cafae10c07b062ed010eee3aad7bd203))


### Tests

* collapse the duplicated guardrail and validator scaffolding ([27dd653](https://github.com/mkutlak/terraform-provider-defectdojo/commit/27dd653f9290b6339faa6b28b1acb45c243e582f))
* fold the standalone acceptance cases into their resource tests ([ec7fb02](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ec7fb02da7deba99c46b4fccbc19824398f189fd))
* guard the round-trip contract between schema flags and the API ([ec54666](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ec546663e61b5a3d15e0f799cef96c1fa34197a1))
* **jira_instance:** re-enable the acceptance tests ([680e0bc](https://github.com/mkutlak/terraform-provider-defectdojo/commit/680e0bcad696b00a08dff89ce6abdc3e72e4a022))

## [1.0.1](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v1.0.0...v1.0.1) (2026-08-06)


### Bug Fixes

* surface conversion diagnostics instead of discarding them ([429b2c1](https://github.com/mkutlak/terraform-provider-defectdojo/commit/429b2c157254c942cd4b026e35185b7d15f13228)), closes [#23](https://github.com/mkutlak/terraform-provider-defectdojo/issues/23)
* **test:** accept date-only target_start/target_end and keep the literal ([7820fbb](https://github.com/mkutlak/terraform-provider-defectdojo/commit/7820fbbc84e8b95345c619e1ded5f51211c5f5c2)), closes [#23](https://github.com/mkutlak/terraform-provider-defectdojo/issues/23)
* **test:** close seeded response bodies to satisfy bodyclose ([50a2cd1](https://github.com/mkutlak/terraform-provider-defectdojo/commit/50a2cd1632683ca58e1bedb6aa08edde8a7f0aa3))


### Miscellaneous

* **deps:** bump actions/setup-go from 6 to 7 ([18568a0](https://github.com/mkutlak/terraform-provider-defectdojo/commit/18568a0591e69a7d39f4caf558313e283803efae))
* **deps:** bump github.com/oapi-codegen/oapi-codegen/v2 ([eec7151](https://github.com/mkutlak/terraform-provider-defectdojo/commit/eec71513aa0d7e1ca9eec01f7a57aa679f802ecc))
* **deps:** update Go module dependencies ([d81cd1e](https://github.com/mkutlak/terraform-provider-defectdojo/commit/d81cd1ef524ce4ea7c0b61bdb7596e0a727656d5))


### CI/CD

* add static analysis with golangci-lint, vet and gofmt gates ([a78c96a](https://github.com/mkutlak/terraform-provider-defectdojo/commit/a78c96ab295b771eef874eda04083c2b58ecd563))
* bump golangci-lint-action and setup-opentofu to node24 runtimes ([f533d09](https://github.com/mkutlak/terraform-provider-defectdojo/commit/f533d09e553dbda4fe080e4be48ca0f109482faa))
* split unit tests from acceptance and gate the matrix on release PRs ([92e3065](https://github.com/mkutlak/terraform-provider-defectdojo/commit/92e30656d15f8ff9aaf1439b5f236f79e8d769ac))
* verify the provider against OpenTofu ([7b9e837](https://github.com/mkutlak/terraform-provider-defectdojo/commit/7b9e8376e72fe553bdeaed945831c92ce3d97686))


### Tests

* add generic acceptance harness with derived resource registry ([d3e837d](https://github.com/mkutlak/terraform-provider-defectdojo/commit/d3e837de0f692efdba97b4533f3f2912dac81e1c))
* cover date handling and diagnostics delivery ([0d14240](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0d14240ca3cd779ff8d7ca09b5e3efd3a6088648)), closes [#23](https://github.com/mkutlak/terraform-provider-defectdojo/issues/23)
* cover the 16 data sources that had no acceptance tests ([2e40c86](https://github.com/mkutlak/terraform-provider-defectdojo/commit/2e40c86701492b19c08daae30e16d704ea88d7aa))
* enforce CheckDestroy on every acceptance TestCase ([dc55692](https://github.com/mkutlak/terraform-provider-defectdojo/commit/dc55692114340a03b86669963a4bcab8ba090971))
* migrate assertions to statecheck and knownvalue ([3b04e69](https://github.com/mkutlak/terraform-provider-defectdojo/commit/3b04e690a5f58c49cdce71d7c3d6070d23df807b))
* validate every resource and data source schema ([a075560](https://github.com/mkutlak/terraform-provider-defectdojo/commit/a075560a5c4d73cabbc61cbaf847d72fda9e7e39))
* verify destroy for all 28 resources ([c279d6a](https://github.com/mkutlak/terraform-provider-defectdojo/commit/c279d6aa8a732fc83a39e0f8b9aca0ed4139399d))

## [1.0.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.5.1...v1.0.0) (2026-07-17)


### ⚠ BREAKING CHANGES

* remove resources dropped by the DefectDojo 3.x API

### Features

* add authorized_users attribute to product and product_type ([3b39b27](https://github.com/mkutlak/terraform-provider-defectdojo/commit/3b39b270850f197009bfd78026e154ea20a6d99a))
* **announcement:** add announcement resource and data source ([ed1ba8a](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ed1ba8a314b3dcd9bd927de90765c67e748b9687))
* **configuration_permission:** add configuration permission data source ([61cb729](https://github.com/mkutlak/terraform-provider-defectdojo/commit/61cb72974351b7629526f03e06544bb517687dd0))
* **endpoint:** deprecate endpoint data source in favor of url/location ([a14694f](https://github.com/mkutlak/terraform-provider-defectdojo/commit/a14694f7b12223b501e3c165c44887d1f788ecd3))
* **engine:** support singleton adoption and defined-element sets ([4155d83](https://github.com/mkutlak/terraform-provider-defectdojo/commit/4155d834c2b5bb747394a7ebb77c46b2b40d9a89))
* **location_product:** add location product join resource and data source ([4d13d86](https://github.com/mkutlak/terraform-provider-defectdojo/commit/4d13d864595ecccdd0a94f92f23e7ee5c6ddd6f6))
* **location:** add read-only location data source ([344b0b7](https://github.com/mkutlak/terraform-provider-defectdojo/commit/344b0b712a4ec2ce851693f6062de576b150e7e7))
* **metadata:** add metadata resource and data source ([9019d63](https://github.com/mkutlak/terraform-provider-defectdojo/commit/9019d635304f8be7147f2d9013ec26814b7a7bea))
* **notifications:** add notifications resource and data source ([3e10bc5](https://github.com/mkutlak/terraform-provider-defectdojo/commit/3e10bc5910cab6144d68d7a085d0ec79d4f31a2d))
* regenerate ddclient from DefectDojo 3.1.101 spec ([5517f52](https://github.com/mkutlak/terraform-provider-defectdojo/commit/5517f524954fd9f59b88a91f79d7924552c6a126))
* remove resources dropped by the DefectDojo 3.x API ([a941841](https://github.com/mkutlak/terraform-provider-defectdojo/commit/a941841990cc79a846be0ac981eab786777bda62))
* **system_settings:** add singleton system settings resource ([d34750e](https://github.com/mkutlak/terraform-provider-defectdojo/commit/d34750e3f184310919fadd4e1d17a8dc45081365))
* **test_type:** add test type data source ([b25e203](https://github.com/mkutlak/terraform-provider-defectdojo/commit/b25e2030bf3f28aecae08c39f89092be80c365f5))
* **url:** add url resource and data source ([0cd76e7](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0cd76e7c8085818e2c2856929c19f6eba04c83c0))
* **user_contact_info:** expose deduplication_execution_mode ([8c60510](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8c60510dec77140d356370c6f802ef45a8736ebe))
* **user_profile:** add current-identity data source ([932bf02](https://github.com/mkutlak/terraform-provider-defectdojo/commit/932bf02f5e9cd7f95a5d552b0ffb7f92caae2853))
* **user:** expose is_staff ([68ecd39](https://github.com/mkutlak/terraform-provider-defectdojo/commit/68ecd39a626f2f85d8756a8a6ac7332fbe398e2e))


### Bug Fixes

* expose authorized_users in product and product_type data sources ([6b59d5d](https://github.com/mkutlak/terraform-provider-defectdojo/commit/6b59d5d99f10868f8a7e12a7aac511557e2c5b29))
* **metadata:** drop location and endpoint parents unsupported by the 3.1 API ([25fb368](https://github.com/mkutlak/terraform-provider-defectdojo/commit/25fb368f22079a69ecc1640d2839b7ffd2fcc53e))
* **notifications:** parse asymmetric scan_added_empty wire format ([f1a2d3c](https://github.com/mkutlak/terraform-provider-defectdojo/commit/f1a2d3c199ffd8ae2fdce8a14ef0df07576ea4c4))
* **url:** mark server-filled optional attributes as computed ([be3f117](https://github.com/mkutlak/terraform-provider-defectdojo/commit/be3f117bb756d9ce7310e1f8fb8ad324c5f877db))
* **url:** reject leading-slash paths that DefectDojo would normalize ([2502cb5](https://github.com/mkutlak/terraform-provider-defectdojo/commit/2502cb5fa3c21fcc137f053e85adb2a7834fe16e))


### Miscellaneous

* default DD_VERSION to 3.1.101 ([0f7f591](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0f7f59137549607bc28af0898b0ebe189d0089e4))
* **deps:** update Go module dependencies ([9e2fd47](https://github.com/mkutlak/terraform-provider-defectdojo/commit/9e2fd47d02c4b350dc270fb964191b209d10256e))
* fix gofmt alignment drift ([5d78712](https://github.com/mkutlak/terraform-provider-defectdojo/commit/5d78712fda07b9f6634b1ed7f98dd9a63274a334))
* fix lint findings in notifications raw client calls and audit test ([227fef5](https://github.com/mkutlak/terraform-provider-defectdojo/commit/227fef5c2eb5c4d8616b2351aa068fff4889b1d2))
* keep collected openapi specs untracked ([56f89e5](https://github.com/mkutlak/terraform-provider-defectdojo/commit/56f89e5ba28fa80567062ffe3e7740663719c92f))
* release the DefectDojo 3.x line as 1.0.0 ([87c9440](https://github.com/mkutlak/terraform-provider-defectdojo/commit/87c9440d82fb1f56363686059adebcf973941a89))


### Documentation

* drop DefectDojo 2.x from supported versions ([3e15b40](https://github.com/mkutlak/terraform-provider-defectdojo/commit/3e15b4045cff4989f20c28f7e74df4cbe8bad2d6))
* remove never-implemented resources from README tables ([d3cc085](https://github.com/mkutlak/terraform-provider-defectdojo/commit/d3cc085038f02d5a8f859e185b86c1d48685c52e))
* update resource counts and exclusion notes for 3.1 coverage ([5fdd3b9](https://github.com/mkutlak/terraform-provider-defectdojo/commit/5fdd3b9194ddacdb95fddb00807228bee88488be))


### CI/CD

* checkout before setup-go so module caching works ([2b75090](https://github.com/mkutlak/terraform-provider-defectdojo/commit/2b75090251371f5a9b24ca37642ea0358f35479f))
* drop DefectDojo 2.x from compat workflow defaults ([621d7be](https://github.com/mkutlak/terraform-provider-defectdojo/commit/621d7becc5c8b5b25dee7d7bed7d472a61ab63d2))


### Tests

* add ddField coverage audit against generated client structs ([c806925](https://github.com/mkutlak/terraform-provider-defectdojo/commit/c80692536931c012a20b0d0c72dadde0a985d25c))
* add ddField reflection audit for generated client drift ([dcb7cf6](https://github.com/mkutlak/terraform-provider-defectdojo/commit/dcb7cf6b5c06aca4c74958091d8845f723b06ece))
* add endpoint data source acceptance test seeded via locations API ([8bee35b](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8bee35bd9a7437d2864befea59f19c5e7d5bd87a))
* **configuration_permission:** look up an existing configuration codename ([bfaf009](https://github.com/mkutlak/terraform-provider-defectdojo/commit/bfaf00932514246061166d0effba3d1626563e53))
* extract shared location seeding helpers for acceptance tests ([d45f77a](https://github.com/mkutlak/terraform-provider-defectdojo/commit/d45f77a75903739109b97dd37b489aa5ac103237))
* **metadata:** match framework validator error in conflict test ([4db9a9b](https://github.com/mkutlak/terraform-provider-defectdojo/commit/4db9a9b0bc7ab34494d43f88cdb20a7632f22c7d))
* replace deprecated parser.ParseDir in ddField audit ([7f16d8a](https://github.com/mkutlak/terraform-provider-defectdojo/commit/7f16d8a29390607b4934b5d05a76632d9a559628))

## [0.5.1](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.5.0...v0.5.1) (2026-07-16)


### Bug Fixes

* pin PGDATA for postgres 18 volume compatibility ([dd2d7f3](https://github.com/mkutlak/terraform-provider-defectdojo/commit/dd2d7f31361c9a56e6a8b82a8aebf8dfe4401555))


### Miscellaneous

* **deps:** bump actions/checkout from 6 to 7 ([0c311aa](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0c311aa7faf788c8133ddb12e15bbc9945674937))
* **deps:** bump github.com/oapi-codegen/runtime from 1.4.1 to 1.4.2 ([b61fdce](https://github.com/mkutlak/terraform-provider-defectdojo/commit/b61fdcea75542f36857ece36712b2f03e3fd7a1d))

## [0.5.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.4.0...v0.5.0) (2026-06-10)


### Features

* target DefectDojo 2.58.4 and Go 1.26.4 ([791d352](https://github.com/mkutlak/terraform-provider-defectdojo/commit/791d352a314404dec107da26609a8c432de5dd10))


### Miscellaneous

* **deps:** bump github.com/hashicorp/terraform-plugin-docs ([2f9e35c](https://github.com/mkutlak/terraform-provider-defectdojo/commit/2f9e35cac1ca1448516debb9ea90ffff5a20bde0))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([4333056](https://github.com/mkutlak/terraform-provider-defectdojo/commit/43330567dd66efaa6b6b808ba8d818eb4000ed06))
* **deps:** bump github.com/hashicorp/terraform-plugin-testing ([9bed886](https://github.com/mkutlak/terraform-provider-defectdojo/commit/9bed8861be4ce71781d090bd49e23c8de5b99293))
* **deps:** bump github.com/oapi-codegen/oapi-codegen/v2 ([51575e2](https://github.com/mkutlak/terraform-provider-defectdojo/commit/51575e20ecd3e86b156c8824194d5715c692c203))
* **deps:** bump github.com/oapi-codegen/runtime from 1.3.0 to 1.3.1 ([01c1823](https://github.com/mkutlak/terraform-provider-defectdojo/commit/01c18232553c9d6144eab295c9bf3b348771681a))
* **deps:** bump github.com/oapi-codegen/runtime from 1.3.1 to 1.4.0 ([de6edd8](https://github.com/mkutlak/terraform-provider-defectdojo/commit/de6edd83cb26fe4132a6569f12c0c27117ca3db8))
* **deps:** bump github.com/oapi-codegen/runtime from 1.4.0 to 1.4.1 ([26b2096](https://github.com/mkutlak/terraform-provider-defectdojo/commit/26b2096484584aa9693de1ddbe987b082b30c7d4))
* **deps:** bump googleapis/release-please-action from 4 to 5 ([629360c](https://github.com/mkutlak/terraform-provider-defectdojo/commit/629360c4d3dc99554b0d6fb9826ce1a40a5cd604))

## [0.4.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.3.1...v0.4.0) (2026-03-24)


### Features

* support name-based lookup for multiple data sources ([8952ac9](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8952ac96ec01513c91070c467385871be2d420f2))

## [0.3.1](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.3.0...v0.3.1) (2026-03-23)


### Bug Fixes

* set default empty string for engagement description ([db72468](https://github.com/mkutlak/terraform-provider-defectdojo/commit/db72468bbd884d3d65fe5e2cca68ab9e9acdf945))


### Tests

* verify empty description in engagement acceptance test ([99a9c52](https://github.com/mkutlak/terraform-provider-defectdojo/commit/99a9c527ba3e7c0ead7c887e255cb6ce48e6e9f8))

## [0.3.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.2.0...v0.3.0) (2026-03-23)


### Features

* improve provider error reporting and schema definitions ([04ff8c4](https://github.com/mkutlak/terraform-provider-defectdojo/commit/04ff8c4459a0190cd9473a13bcd7895ae8b1a8b7))


### Miscellaneous

* add multi-version DefectDojo compatibility testing and automation ([0c2eabb](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0c2eabb8a5a5ae9613cf3d5f3beb5a6f706f5d90))
* update dependencies and add generate-docs target ([97744b2](https://github.com/mkutlak/terraform-provider-defectdojo/commit/97744b2f0e02eb69b3802790062ab3b3fd6c4eba))

## [0.2.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.1.2...v0.2.0) (2026-03-18)


### Features

* add comprehensive set of resources and data sources ([20bd4d2](https://github.com/mkutlak/terraform-provider-defectdojo/commit/20bd4d28b0401055f5a2c5ff9bbb30e66b0cc6cf))
* implement 0.1.0-rc1 changes and refactor to a new namespace ([3da79fe](https://github.com/mkutlak/terraform-provider-defectdojo/commit/3da79fe3172989ccad8a31f6b87d5ba3dfb3e7df))


### Bug Fixes

* enable parallel execution for acceptance tests and update manifest to v0.1.1 ([191079f](https://github.com/mkutlak/terraform-provider-defectdojo/commit/191079f7b3c750bdde746e5eb96058e39b3fa7c3))
* improve api error handling and resource data population logic ([ffe884f](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ffe884f0d4b67d2d0d4e25da5020815f2a5bc20c))


### Miscellaneous

* **master:** release 0.1.0 ([302f2cd](https://github.com/mkutlak/terraform-provider-defectdojo/commit/302f2cd095053b33ab1f50b7de85fcd47680e38a))
* **master:** release 0.1.0 ([f4b69ea](https://github.com/mkutlak/terraform-provider-defectdojo/commit/f4b69eab5f12b7c57f6b010aefba75c2d9f6336c))
* **master:** release 0.1.2 ([05e51fe](https://github.com/mkutlak/terraform-provider-defectdojo/commit/05e51fe7f38ef8be3b132db3cda8db9626952b79))
* **master:** release 0.1.2 ([503c3d8](https://github.com/mkutlak/terraform-provider-defectdojo/commit/503c3d8e337aeb423352f33e8335a8517b72ad38))
* update CI workflows, add CLAUDE.md, and update .gitignore ([8df7195](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8df719594049f854041d06143c4d3c33004180a0))
* update dependencies in go.mod and go.sum ([6df6327](https://github.com/mkutlak/terraform-provider-defectdojo/commit/6df632715d65bc39d3d483b9e52e365fcabdef67))


### CI/CD

* add setup-terraform to test workflow ([ab474c7](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ab474c7c979f9af0e36e228a0a108d701372734b))
* Add token to release-please workflow ([0656107](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0656107a7e081eb4670bebb63e1888b609dde84c))
* Fix version from 0.1.0-rc1 to 0.0.0 ([4f3b2dc](https://github.com/mkutlak/terraform-provider-defectdojo/commit/4f3b2dc20a5c327b04712c0cf18c90d2aa94c4af))
* setup release-please and improve local testing environment ([416f653](https://github.com/mkutlak/terraform-provider-defectdojo/commit/416f653b22a763742565be279cbc284f0a2dde96))

## [0.1.2](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.1.1...v0.1.2) (2026-03-18)


### Bug Fixes

* enable parallel execution for acceptance tests and update manifest to v0.1.1 ([191079f](https://github.com/mkutlak/terraform-provider-defectdojo/commit/191079f7b3c750bdde746e5eb96058e39b3fa7c3))
* improve api error handling and resource data population logic ([ffe884f](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ffe884f0d4b67d2d0d4e25da5020815f2a5bc20c))


### Miscellaneous

* update dependencies in go.mod and go.sum ([6df6327](https://github.com/mkutlak/terraform-provider-defectdojo/commit/6df632715d65bc39d3d483b9e52e365fcabdef67))

## [0.1.0](https://github.com/mkutlak/terraform-provider-defectdojo/compare/v0.0.0...v0.1.0) (2026-03-18)


### Features

* implement 0.1.0-rc1 changes and refactor to a new namespace ([3da79fe](https://github.com/mkutlak/terraform-provider-defectdojo/commit/3da79fe3172989ccad8a31f6b87d5ba3dfb3e7df))


### Miscellaneous

* update CI workflows, add CLAUDE.md, and update .gitignore ([8df7195](https://github.com/mkutlak/terraform-provider-defectdojo/commit/8df719594049f854041d06143c4d3c33004180a0))


### CI/CD

* add setup-terraform to test workflow ([ab474c7](https://github.com/mkutlak/terraform-provider-defectdojo/commit/ab474c7c979f9af0e36e228a0a108d701372734b))
* Add token to release-please workflow ([0656107](https://github.com/mkutlak/terraform-provider-defectdojo/commit/0656107a7e081eb4670bebb63e1888b609dde84c))
* Fix version from 0.1.0-rc1 to 0.0.0 ([4f3b2dc](https://github.com/mkutlak/terraform-provider-defectdojo/commit/4f3b2dc20a5c327b04712c0cf18c90d2aa94c4af))
* setup release-please and improve local testing environment ([416f653](https://github.com/mkutlak/terraform-provider-defectdojo/commit/416f653b22a763742565be279cbc284f0a2dde96))

## 0.0.13

FEATURES
  - Add two attributes to `defectdojo_product_type` resource and data source:
    - `critical_product`
    - `key_product`

## 0.0.12

BUGFIX:
 - A product with no tags specified would cause a provider error from terraform.

## 0.0.11

FEATURES:
 - Add the following fields to `defectdojo_product` resource:
   - `business_criticality`
   - `enable_full_risk_acceptance`
   - `enable_skip_risk_acceptance`
   - `external_audience`
   - `internet_accessible`
   - `lifecycle`
   - `origin`
   - `platform`
   - `prod_numeric_grade`
   - `regulation_ids`
   - `revenue`
   - `user_records`

## 0.0.10

FEATURES:
 - Add `jira_product_configuration` resource.

## 0.0.9

BUGFIX:
 - Fix delete-drift detection in `product` and `product_type` resources. If the object was deleted outside terraform we remove it from the state.

## 0.0.8

BUGFIX:
 - Don't continue processing after encountering an error, that cause a panic.

## 0.0.7

Initial public release

## 0.0.4

FEATURES:
 - Add basic support for Product Type resource and data source

## 0.0.3

FEATURES:
 - First working version.
 - Basic support for Product resource and data source.
