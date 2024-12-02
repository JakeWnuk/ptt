# Security Documentation for Password Transformation Tool (PTT)

## Notices
- As of `v0.4.1`, the HCRE library is included in the project. This library is
  used to support complex rule features and has been locked to the following
  versions:
  - launchpad.net/hcre v0.0.0-20241130145909-c832018180b1 h1:lfPqGETHlSypBMeJtjMAFTthkaE/Wkxgu1vzYpwKdEI=
  - launchpad.net/hcre v0.0.0-20241130145909-c832018180b1/go.mod h1:Dq78e8vypvdrOQt+VImkJcRq/6GHM1XGvO9/T1nr18M=

## Notes:
- `pkg/utils/utils.go:9` use of `crypto/rand` over `math/rand` is not needed in this module.

## Last SAST Scan:
- `v0.4.1`
    - Included `hcre` library
