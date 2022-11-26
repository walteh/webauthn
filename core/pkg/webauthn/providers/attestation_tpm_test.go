package providers

import (
	"crypto/sha256"
	"testing"

	"nugg-webauthn/core/pkg/errors"
	"nugg-webauthn/core/pkg/hex"
	"nugg-webauthn/core/pkg/webauthn/credential"
	"nugg-webauthn/core/pkg/webauthn/googletpm"
	"nugg-webauthn/core/pkg/webauthn/types"

	"github.com/stretchr/testify/assert"
)

var provider = TpmAttestationProvider{}

func TestTPMAttestationVerificationSuccess(t *testing.T) {
	for i := range testAttestationTPMResponses {
		t.Run("TPM Positive tests", func(t *testing.T) {
			pcc := testAttestationTPMResponses[i]
			clientDataHash := sha256.Sum256([]byte(pcc.UTF8ClientDataJSON))

			att, err := credential.ParseAttestationInput(pcc)
			assert.NoError(t, err)

			attestationKey := provider.ID()
			_, _, _, err = provider.Handler(*att, clientDataHash[:])
			if err != nil {
				t.Fatalf("Not valid: %+v", err)
			}
			assert.Equal(t, "tpm", attestationKey)
		})
	}
}

var testAttestationTPMResponses = []types.AttestationInput{
	// TPM attestation with ECC P256
	{
		CredentialID:       hex.MustBase64ToHash("hsS2ywFz_LWf9-lC35vC9uJTVD3ZCVdweZvESUbjXnQ"),
		CredentialType:     "public-key",
		UTF8ClientDataJSON: "{\"type\":\"webauthn.create\",\"challenge\":\"uzn9u0Tx-LBdtGgERsbkHRBjiUt5i2rvm2BBTZrWqEo\",\"origin\":\"https://webauthn.io\",\"crossOrigin\":false}",
		AttestationObject:  hex.MustBase64ToHash("o2NmbXRjdHBtZ2F0dFN0bXSmY2FsZzn__mNzaWdZAQCqAcGoi2IFXCF5xxokjR5yOAwK_11iCOqt8hCkpHE9rW602J3KjhcRQzoFf1UxZvadwmYcHHMxDQDmVuOhH-yW-DfARVT7O3MzlhhzrGTNO_-jhGFsGeEdz0RgNsviDdaVP5lNsV6Pe4bMhgBv1aTkk0zx1T8sxK8B7gKT6x80RIWg89_aYY4gHR4n65SRDp2gOGI2IHDvqTwidyeaAHVPbDrF8iDbQ88O-GH_fheAtFtgjbIq-XQbwVdzQhYdWyL0XVUwGLSSuABuB4seRPkyZCKoOU6VuuQzfWNpH2Nl05ybdXi27HysUexgfPxihB3PbR8LJdi1j04tRg3JvBUvY3ZlcmMyLjBjeDVjglkFuzCCBbcwggOfoAMCAQICEGEZiaSlAkKpqaQOKDYmWPkwDQYJKoZIhvcNAQELBQAwQTE_MD0GA1UEAxM2RVVTLU5UQy1LRVlJRC1FNEE4NjY2RjhGNEM2RDlDMzkzMkE5NDg4NDc3ODBBNjgxMEM0MjEzMB4XDTIyMDExMjIyMTUxOFoXDTI3MDYxMDE4NTQzNlowADCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKo-7DHdiipZTzfA9fpTaIMVK887zM0nXAVIvU0kmGAsPpTYbf7dn1DAl6BhcDkXs2WrwYP02K8RxXWOF4jf7esMAIkr65zPWqLys8WRNM60d7g9GOADwbN8qrY0hepSsaJwjhswbNJI6L8vJwnnrQ6UWVCm3xHqn8CB2iSWNSUnshgTQTkJ1ZEdToeD51sFXUE0fSxXjyIiSAAD4tCIZkmHFVqchzfqUgiiM_mbbKzUnxEZ6c6r39ccHzbm4Ir-u62repQnVXKTpzFBbJ-Eg15REvw6xuYaGtpItk27AXVcEodfAylf7pgQPfExWkoMZfb8faqbQAj5x29mBJvlzj0CAwEAAaOCAeowggHmMA4GA1UdDwEB_wQEAwIHgDAMBgNVHRMBAf8EAjAAMG0GA1UdIAEB_wRjMGEwXwYJKwYBBAGCNxUfMFIwUAYIKwYBBQUHAgIwRB5CAFQAQwBQAEEAIAAgAFQAcgB1AHMAdABlAGQAIAAgAFAAbABhAHQAZgBvAHIAbQAgACAASQBkAGUAbgB0AGkAdAB5MBAGA1UdJQQJMAcGBWeBBQgDMFAGA1UdEQEB_wRGMESkQjBAMT4wEAYFZ4EFAgIMB05QQ1Q3NXgwFAYFZ4EFAgEMC2lkOjRFNTQ0MzAwMBQGBWeBBQIDDAtpZDowMDA3MDAwMjAfBgNVHSMEGDAWgBQ3yjAtSXrnaSNOtzy1PEXxOO1ZUDAdBgNVHQ4EFgQU1ml3H5Tzrs0Nev69tFNhPZnhaV0wgbIGCCsGAQUFBwEBBIGlMIGiMIGfBggrBgEFBQcwAoaBkmh0dHA6Ly9hemNzcHJvZGV1c2Fpa3B1Ymxpc2guYmxvYi5jb3JlLndpbmRvd3MubmV0L2V1cy1udGMta2V5aWQtZTRhODY2NmY4ZjRjNmQ5YzM5MzJhOTQ4ODQ3NzgwYTY4MTBjNDIxMy9lMDFjMjA2Mi1mYmRjLTQwYTUtYTQwZi1jMzc3YzBmNzY1MWMuY2VyMA0GCSqGSIb3DQEBCwUAA4ICAQAz-YGrj0S841gyMZuit-qsKpKNdxbkaEhyB1baexHGcMzC2y1O1kpTrpaH3I80hrIZFtYoA2xKQ1j67uoC6vm1PhsJB6qhs9T7zmWZ1VtleJTYGNZ_bYY2wo65qJHFB5TXkevJUVe2G39kB_W1TKB6g_GSwb4a5e4D_Sjp7b7RZpyIKHT1_UE1H4RXgR9Qi68K4WVaJXJUS6T4PHrRc4PeGUoJLQFUGxYokWIf456G32GwGgvUSX76K77pVv4Y-kT3v5eEJdYxlS4EVT13a17KWd0DdLje0Ae69q_DQSlrHVLUrADvuZMeM8jxyPQvDb7ETKLsSUeHm73KOCGLStcGQ3pB49nt3d9XdWCcUwUrmbBF2G7HsRgTNbj16G6QUcWroQEqNrBG49aO9mMZ0NwSn5d3oNuXSXjLdGBXM1ukLZ-GNrZDYw5KXU102_5VpHpjIHrZh0dXg3Q9eucKe6EkFbH65-O5VaQWUnR5WJpt6-fl_l0iHqHnKXbgL6tjeerCqZWDvFsOak05R-hosAoQs_Ni0EsgZqHwR_VlG86fsSwCVU3_sDKTNs_Je08ewJ_bbMB5Tq6k1Sxs8Aw8R96EwjQLp3z-Zva1myU-KerYYVDl5BdvgPqbD8Xmst-z6vrP3CJbtr8jgqVS7RWy_cJOA8KCZ6IS_75QT7Gblq6UGFkG7zCCBuswggTToAMCAQICEzMAAAbTtnznKsOrB-gAAAAABtMwDQYJKoZIhvcNAQELBQAwgYwxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpXYXNoaW5ndG9uMRAwDgYDVQQHEwdSZWRtb25kMR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24xNjA0BgNVBAMTLU1pY3Jvc29mdCBUUE0gUm9vdCBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkgMjAxNDAeFw0yMTA2MTAxODU0MzZaFw0yNzA2MTAxODU0MzZaMEExPzA9BgNVBAMTNkVVUy1OVEMtS0VZSUQtRTRBODY2NkY4RjRDNkQ5QzM5MzJBOTQ4ODQ3NzgwQTY4MTBDNDIxMzCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAJA7GLwHWWbn2H8DRppxQfre4zll1sgE3Wxt9DTYWt5-v-xKwCQb6z_7F1py7LMe58qLqglAgVhS6nEvN2puZ1GzejdsFFxz2gyEfH1y-X3RGp0dxS6UKwEtmksaMEKIRQn2GgKdUkiuvkaxaoznuExoTPyu0aXk6yFsX5KEDu9UZCgt66bRy6m3KIRnn1VK2frZfqGYi8C8x9Q69oGG316tUwAIm3ypDtv3pREXsDLYE1U5Irdv32hzJ4CqqPyau-qJS18b8CsjvgOppwXRSwpOmU7S3xqo-F7h1eeFw2tgHc7PEPt8MSSKeba8Fz6QyiLhgFr8jFUvKRzk4B41HFUMqXYawbhAtfIBiGGsGrrdNKb7MxISnH1E6yLVCQGGhXiN9U7V0h8Gn56eKzopGlubw7yMmgu8Cu2wBX_a_jFmIBHnn8YgwcRm6NvT96KclDHnFqPVm3On12bG31F7EYkIRGLbaTT6avEu9rL6AJn7Xr245Sa6dC_OSMRKqLSufxp6O6f2TH2g4kvT0Go9SeyM2_acBjIiQ0rFeBOm49H4E4VcJepf79FkljovD68imeZ5MXjxepcCzS138374Jeh7k28JePwJnjDxS8n9Dr6xOU3_wxS1gN5cW6cXSoiPGe0JM4CEyAcUtKrvpUWoTajxxnylZuvS8ou2thfH2PQlAgMBAAGjggGOMIIBijAOBgNVHQ8BAf8EBAMCAoQwGwYDVR0lBBQwEgYJKwYBBAGCNxUkBgVngQUIAzAWBgNVHSAEDzANMAsGCSsGAQQBgjcVHzASBgNVHRMBAf8ECDAGAQH_AgEAMB0GA1UdDgQWBBQ3yjAtSXrnaSNOtzy1PEXxOO1ZUDAfBgNVHSMEGDAWgBR6jArOL0hiF-KU0a5VwVLscXSkVjBwBgNVHR8EaTBnMGWgY6Bhhl9odHRwOi8vd3d3Lm1pY3Jvc29mdC5jb20vcGtpb3BzL2NybC9NaWNyb3NvZnQlMjBUUE0lMjBSb290JTIwQ2VydGlmaWNhdGUlMjBBdXRob3JpdHklMjAyMDE0LmNybDB9BggrBgEFBQcBAQRxMG8wbQYIKwYBBQUHMAKGYWh0dHA6Ly93d3cubWljcm9zb2Z0LmNvbS9wa2lvcHMvY2VydHMvTWljcm9zb2Z0JTIwVFBNJTIwUm9vdCUyMENlcnRpZmljYXRlJTIwQXV0aG9yaXR5JTIwMjAxNC5jcnQwDQYJKoZIhvcNAQELBQADggIBAFZTSitCISvll6i6rPUPd8Wt2mogRw6I_c-dWQzdc9-SY9iaIGXqVSPKKOlAYU2ju7nvN6AvrIba6sngHeU0AUTeg1UZ5-bDFOWdSgPaGyH_EN_l-vbV6SJPzOmZHJOHfw2WT8hjlFaTaKYRXxzFH7PUR4nxGRbWtdIGgQhUlWg5oo_FO4bvLKfssPSONn684qkAVierq-ly1WeqJzOYhd4EylgVJ9NL3YUhg8dYcHAieptDzF7OcDqffbuZLZUx6xcyibhWQcntAh7a3xPwqXxENsHhme_bqw_kqa-NVk-Wz4zdoiNNLRvUmCSL1WLc4JPsFJ08Ekn1kW7f9ZKnie5aw-29jEf6KIBt4lGDD3tXTfaOVvWcDbu92jMOO1dhEIj63AwQiDJgZhqnrpjlyWU_X0IVQlaPBg80AE0Y3sw1oMrY0XwdeQUjSpH6e5fTYKrNB6NMT1jXGjKIzVg8XbPWlnebP2wEhq8rYiDR31b9B9Sw_naK7Xb-Cqi-VQdUtknSjeljusrBpxGUx-EIJci0-dzeXRT5_376vyKSuYxA1Xd2jd4EknJLIAVLT3rb10DCuKGLDgafbsfTBxVoEa9hSjYOZUr_m3WV6t6I9WPYjVyhyi7fCEIG4JE7YbM4na4jg5q3DM8ibE8jyufAq0PfJZTJyi7c2Q2N_9NgnCNwZ3B1YkFyZWFYdgAjAAsABAByACCd_8vzbDg65pn7mGjcbcuJ1xU4hL4oA5IsEkFYv60irgAQABAAAwAQACAek7g2C8TeORRoKxuN7HrJ5OinVGuHzEgYODyUsF9D1wAggXPPXn-Pm_4IF0c4XVaJjmHO3EB2KBwdg_L60N0IL9xoY2VydEluZm9Yof9UQ0eAFwAiAAvQNGTLa2wT6u8SKDDdwkgaq5Cmh6jcD_6ULvM9ZmvdbwAUtMInD3WtGSdWHPWijMrW_TfYo-gAAAABPuBems3Sywu4aQsGAe85iOosjtXIACIAC5FPRiZSJzjYMNnAz9zFtM62o57FJwv8F5gNEcioqhHwACIACyVXxq1wZhDsqTqdYr7vQUUJ3vwWVrlN0ZQv5HFnHqWdaGF1dGhEYXRhWKR0puqSE8mcL3SyJJKzIM9AJiqUwalQoDl_KSULYIQe8EUAAAAACJhwWMrcS4G24TDeUNy-lgAghsS2ywFz_LWf9-lC35vC9uJTVD3ZCVdweZvESUbjXnSlAQIDJiABIVggHpO4NgvE3jkUaCsbjex6yeTop1Rrh8xIGDg8lLBfQ9ciWCCBc89ef4-b_ggXRzhdVomOYc7cQHYoHB2D8vrQ3Qgv3A"),
	},
	// TPM attestation with RSA SHA1
	{
		CredentialID:       hex.MustBase64ToHash("UJDoUJoGiDQF_EEZ3G_z9Lfq16_KFaXtMTjwTUrrRlc"),
		CredentialType:     "public-key",
		UTF8ClientDataJSON: "{\"origin\":\"https://localhost:44329\",\"challenge\":\"9JyUfJkg8PqoKZuD7FHzOE9dbyculC9urGTpGqBnEwnhKmni4rGRXxm3-ZBHK8x6riJQqIpC8qEa-T0qIFTKTQ\",\"type\":\"webauthn.create\"}",
		AttestationObject:  hex.MustBase64ToHash("o2NmbXRjdHBtZ2F0dFN0bXSmY2FsZzn__mNzaWdZAQBIwu9LPAl-LgxlRzPlvn7L-0yuMnFFn1XALxXtGnmC5-oMIIqfUJWFbgBbkN2l2zPsqOCRT5GQU8ucKNI6HrlbuDAUIq7wjcxG5TzgQt3YtGMWtgEcrZn2ecUlQFKjY67_wZIuHLy443Ki1SjErNPrMrkIPe9lyFhIalMgrWLCol40gYIVr_9xLfgyX55c7XiB-XbUKhDLUv5uPA3CSAiWeWwWx26K2BTV85vHsaG6f2YFTfcQTFs1cTSwMm7A9C2SiQ7N01ENwM1urVxlCvuEsBgiXapR70Oyq_cfiENYY0ti7_w2fvikmfv0z0O1cJOAyUlYWjnWhT707chrVmkFY3ZlcmMyLjBjeDVjglkEXzCCBFswggNDoAMCAQICDwRsOt2imXnV5Z4BftcqfzANBgkqhkiG9w0BAQsFADBBMT8wPQYDVQQDEzZOQ1UtTlRDLUtFWUlELTM2MTA0Q0U0MEJCQ0MxRjQwRDg0QTRCQkQ1MEJFOTkwMjREOTU3RDQwHhcNMTgwMjAxMDAwMDAwWhcNMjUwMTMxMjM1OTU5WjAAMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmw-4ficURR_sgVfW7cs1iRoDGdxjBpCczF233ba_5WTP-RrsYZPlzWgSN9WXptuywzjZoDlbid7NlduSR1ZFsds4bW71LyKDL62eyqaiAc645gocXAyxdDIDJAeo-3N9Dm4vsw-Gy_0sd2v1UEkBhWjuE1gL5hcaB9EtXSDvHPwmrf0eYn_4cWu9AxqSxpn79JIPYEOUrURr2H8zyG4_P0j1a3MVBmtAymhpXBn9ila-bW7K_k0JYXBh5yAYZDsmHgFsXbUauDWdja3HYzkep9jXkFcegXOMjPr_QSqWRjawEvzoprnJ-QqoWNbaRhuD-UnfgCNbwseU8kZ0aQNjBQIDAQABo4IBjzCCAYswDgYDVR0PAQH_BAQDAgeAMAwGA1UdEwEB_wQCMAAwUwYDVR0gAQH_BEkwRzBFBgkrBgEEAYI3FR8wODA2BggrBgEFBQcCAjAqEyhGQUtFIEZJRE8gVENQQSBUcnVzdGVkIFBsYXRmb3JtIElkZW50aXR5MBAGA1UdJQQJMAcGBWeBBQgDMEoGA1UdEQEB_wRAMD6kPDA6MTgwDgYFZ4EFAgMMBWlkOjEzMBAGBWeBBQICDAdOUENUNnh4MBQGBWeBBQIBDAtpZDpGRkZGRjFEMDAfBgNVHSMEGDAWoBRRfyLI5lOlfNVM3TBYfjD_ZzaMXTAdBgNVHQ4EFgQUO6SUmiOhCHVZcq-88acg2uQkQz8weAYIKwYBBQUHAQEEbDBqMGgGCCsGAQUFBzAChlxodHRwczovL2ZpZG9hbGxpYW5jZS5jby5uei90cG1wa2kvTkNVLU5UQy1LRVlJRC0zNjEwNENFNDBCQkNDMUY0MEQ4NEE0QkJENTBCRTk5MDI0RDk1N0Q0LmNydDANBgkqhkiG9w0BAQsFAAOCAQEAIIyVBkck_SD2nbj4KOwUI6cYZHrjwrcULoEiOSXn9TjTIiB5MdBMvqqNyAXiyWoWd1GEc_MI3mKOzu4g5UTVQQqfiOTrqfuZrpoU0tAeojKnZLj2wYj5GpyOfEkPK3m9qVaDxiYrh6aS8a3w_Iog878EiIaoVALbBt5uAfh0TAHHwSdxHtU8DRJrC43yIqcP9byRqssJmgSNcpMAjw_hcKJxDMD2UurvsMasqyWvK533yNA0-VwXvk3HI0ItSOw_g352D-qOTHI82lJIjc3yKoaNeYKn7RzgcLAF7AesTiiJReY2kU_vLyf-wH54-08T3oyBBJpBCHc1y_Lt5d2qWFkGCDCCBgQwggPsoAMCAQICENBTpEeEh5lpTgeR7VT9oQcwDQYJKoZIhvcNAQELBQAwgb8xCzAJBgNVBAYTAlVTMQswCQYDVQQIDAJNWTESMBAGA1UEBwwJV2FrZWZpZWxkMRYwFAYDVQQKDA1GSURPIEFsbGlhbmNlMQwwCgYDVQQLDANDV0cxNjA0BgNVBAMMLUZJRE8gRmFrZSBUUE0gUm9vdCBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkgMjAxODExMC8GCSqGSIb3DQEJARYiY29uZm9ybWFuY2UtdG9vbHNAZmlkb2FsbGlhbmNlLm9yZzAeFw0xNzAyMDEwMDAwMDBaFw0zNTAxMzEyMzU5NTlaMEExPzA9BgNVBAMTNk5DVS1OVEMtS0VZSUQtMzYxMDRDRTQwQkJDQzFGNDBEODRBNEJCRDUwQkU5OTAyNEQ5NTdENDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANc-c30RpQd-_LCoiLJbXz3t_vqciOIovwjez79_DtVgi8G9Ph-tPL-lC0ueFGBMSPcKd_RDdSFe2QCYQd9e0DtiFxra-uWGa0olI1hHI7bK2GzNAZSTKEbwgqpf8vXMQ-7SPajg6PfxSOLH_Nj2yd6tkNkUSdlGtWfY8XGB3n-q--nt3UHdUQWEtgUoTe5abBXsG7MQSuTNoad3v6vk-tLd0W44ivM6pbFqFUHchx8mGLApCpjlVXrfROaCoc9E91hG9B-WNvekJ0dM6kJ658Hy7yscQ6JdqIEolYojCtWaWNmwcfv--OE1Ax_4Ub24gl3hpB9EOcBCzpb4UFmLYUECAwEAAaOCAXcwggFzMAsGA1UdDwQEAwIBhjAWBgNVHSAEDzANMAsGCSsGAQQBgjcVHzAbBgNVHSUEFDASBgkrBgEEAYI3FSQGBWeBBQgDMBIGA1UdEwEB_wQIMAYBAf8CAQAwHwYDVR0OBBgEFsIUUX8iyOZTpXzVTN0wWH4w_2c2jF0wHwYDVR0jBBgwFqAUXH82LZCtWry6jnXa3jqg7cFOAoswaAYDVR0fBGEwXzBdoFugWYZXaHR0cHM6Ly9maWRvYWxsaWFuY2UuY28ubnovdHBtcGtpL2NybC9GSURPIEZha2UgVFBNIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5IDIwMTguY3JsMG8GCCsGAQUFBwEBBGMwYTBfBggrBgEFBQcwAoZTaHR0cHM6Ly9maWRvYWxsaWFuY2UuY28ubnovdHBtcGtpL0ZJRE8gRmFrZSBUUE0gUm9vdCBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkgMjAxOC5jcnQwDQYJKoZIhvcNAQELBQADggIBAG138t55DF9nPJbvbPQZOypmyTPpNne0A5fh69P1fHZ5qdE2PDz3cf5Tl-8OPI4xQniEFNPcXMb7KlhMM6zCl4GkZtNN4MxygdFjQ1gTZOBDpt7Dwziij0MakmwyC0RYTNtbSyVhHUevgw9rnu13EzqxPyL5JD-UqADh2Y51MS0qy7IOgegLQv-eJzSNUgHxFJreUzz4PU6yzSsTyyYDW-H4ZjAQKienVp8ewZf8oHGWHGQFGa5E9m1P8vxCMZ7pIzeQweCVYrs3q7unu4nzBAIXLPI092kYFUgyz3lIaSB3XEiPBokpupX6Zmgrfphb-XX3tbenH5hkxfumueA5RMHTMu5TVjhJXiV0yM3q5W5xrQHdJlF5nOdJDEE-Kb7nm6xaT1DDpafqBc5vEDMkJmBA4AXHUY7JPGqEEzEenT7k6Wn5IQLZg4qc8Irnj__yM7xUhJWJam47KVbLA4WFu-IKvJrkP5GSglZ9qASOCxBHaOL2UcTAg50uvhUSwur2KSak2vlENdmAijwdAL4LLQWrkFd-9NBwcNwTdfK4ekEHP1l4BwJtkNwW6etUgeA5rkW2JLocXoBq5v7GSk4_CBoKhyiahQGQQ9SZFGeBJhzzkK9yN-yKskcVjjjInSHPl-ZpeOK3sI08sEyTH0gxlTtRoX0MKDsMAHEVToe5o1u9Z3B1YkFyZWFZATYAAQALAAYEcgAgnf_L82w4OuaZ-5ho3G3LidcVOIS-KAOSLBJBWL-tIq4AEAAQCAAAAAAAAQCl9siJwqoHJ2pCwEKyLQ_u6zGcZDKZtA0jtvtn1aPlIe7wFAvQNgjI6KDiQsDPTCVeJj_RA441VbV0Z4oX2b68quDY0Gf4VpF4KWfNPdKH6H4E882m8OnBb10mhaNbPxTmDVDZLQZjh3ubX1Z56FNg6cQmz4bEnHF-7X1l7AcNORhzdzgM7uRXhwo9UsAzpu4Io1OCTsb5DaDnng3f3Y9qDn8OG3MI_5IYtm1qGgmY72nSEiIhhPCk2lvmajN6A4tWgUstc7QtdlKEPBd-ITtGdKYTSwqihaHzBQd8D-d_HDqgcOWECLKo51_YqyaEiuGlv6sPon1LMsEL6PlVw47PaGNlcnRJbmZvWKH_VENHgBcAIgALEeaO1E21Ny4UKW4vhKzHg5h1GIGSHjD8IqBvi3PHlFMAFF6MXAvgUX_Rbc04fmdB2TyLG-mdAAAAAUdwF0hVaXtLxoVgpQFzfvmNNFZV-wAiAAuYlrm-5Jg3251TsEdZ8NV11xd4X5O3q0AFLmammw658QAiAAtuzX-04mcxAHq9kO70Ew3vJCOmCS0UvQzZB2CNCeGXpWhhdXRoRGF0YVkBZ0mWDeWIDoxodDQXD2R2YFuP5K65ooYyx5lc87qDHZdjQQAAAHXyRLZ-U2RP1Z-Qw5YicxfbACBQkOhQmgaINAX8QRncb_P0t-rXr8oVpe0xOPBNSutGV6QBAwM5__4gWQEApfbIicKqBydqQsBCsi0P7usxnGQymbQNI7b7Z9Wj5SHu8BQL0DYIyOig4kLAz0wlXiY_0QOONVW1dGeKF9m-vKrg2NBn-FaReClnzT3Sh-h-BPPNpvDpwW9dJoWjWz8U5g1Q2S0GY4d7m19WeehTYOnEJs-GxJxxfu19ZewHDTkYc3c4DO7kV4cKPVLAM6buCKNTgk7G-Q2g554N392Pag5_DhtzCP-SGLZtahoJmO9p0hIiIYTwpNpb5mozegOLVoFLLXO0LXZShDwXfiE7RnSmE0sKooWh8wUHfA_nfxw6oHDlhAiyqOdf2KsmhIrhpb-rD6J9SzLBC-j5VcOOzyFDAQAB"),
	},
	// TPM attestation with RSA SHA256
	{
		CredentialID:       hex.MustBase64ToHash("h9XMhkVePN1Prq9Ks_VfwIsVZvt-jmSRTEnevTc-KB8"),
		CredentialType:     "public-key",
		UTF8ClientDataJSON: "{\"origin\":\"https://localhost:44329\",\"challenge\":\"gHrAk4pNe2VlB0HLeKclI2P6QEa83PuGeijTHMtpbhY9KlybyhlwF_VzRe7yhabXagWuY6rkDWfvvhNqgh2o7A\",\"type\":\"webauthn.create\"}",
		AttestationObject:  hex.MustBase64ToHash("o2NmbXRjdHBtZ2F0dFN0bXSmY2FsZzkBAGNzaWdZAQA6Gh1Oa3-8vCY8bTrpUHA4zp4UCsbuh36tH09G-qWlvQdoqEQsJJQu1Rz61_mFes9CXE2cxiJV8pEwxtUUTSZQWnamVU1x9bBk07qcHqAuamP_NDAahHhZ9D46q9JklT3aVdhbaZVh0y5b8NZB2eUfKqcUmM0JCxLP9ZfSe7XcVguhQVEduM6Qnl9R1zRh7cquOa8UOEpdXkt1-drsOtrA9c0UJPYzkI8qscCDc-xfzo2xv12tLXjRq395JnynHhjzJIz8Ch2IYQUiMSM6TQDcnvzDEvRgril9NC0aIkHd79omIZNnBjEDfjyqOZbBffjGyvt1Eikz4M0EE8e7N4uRY3ZlcmMyLjBjeDVjglkEXzCCBFswggNDoAMCAQICDwQ_ozlil_l5hh6NlMsLzzANBgkqhkiG9w0BAQsFADBBMT8wPQYDVQQDEzZOQ1UtTlRDLUtFWUlELTM2MTA0Q0U0MEJCQ0MxRjQwRDg0QTRCQkQ1MEJFOTkwMjREOTU3RDQwHhcNMTgwMjAxMDAwMDAwWhcNMjUwMTMxMjM1OTU5WjAAMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAor_6-4WYizZdOQ9Ia_offaIdL2BVGtGDq8jQxo16ymBSOWCP15gZt9QAkqowS3ayqEh48Pg5SdA7F5kcjD_FqKaZDBOqkjvJivdo7FKv7EaUI2al9B7h0pXIRb97jn2z0zPlXz6RV_RmBe3CCljyxrhav7bTkCXEJUnkNgxsWgLGBIW6VSVct0z42xBB6_6mYekWIej5vXLqB8AuzsqnLbU5jOohfJiI5urFso12j6YCWZ_kXK4j8e4IoHUOjWgtHXdb3kP8PvI948hcJpIEpuuLDZDDOCOPI1wAlryGwz_tJLarODZzD1XhG3BMlXi1TG7x1s-AriC3A7B89wuSpwIDAQABo4IBjzCCAYswDgYDVR0PAQH_BAQDAgeAMAwGA1UdEwEB_wQCMAAwUwYDVR0gAQH_BEkwRzBFBgkrBgEEAYI3FR8wODA2BggrBgEFBQcCAjAqEyhGQUtFIEZJRE8gVENQQSBUcnVzdGVkIFBsYXRmb3JtIElkZW50aXR5MBAGA1UdJQQJMAcGBWeBBQgDMEoGA1UdEQEB_wRAMD6kPDA6MTgwDgYFZ4EFAgMMBWlkOjEzMBAGBWeBBQICDAdOUENUNnh4MBQGBWeBBQIBDAtpZDpGRkZGRjFEMDAfBgNVHSMEGDAWoBRRfyLI5lOlfNVM3TBYfjD_ZzaMXTAdBgNVHQ4EFgQUS1ZtGu6ZoewTH3mq04Ytxa4kOQcweAYIKwYBBQUHAQEEbDBqMGgGCCsGAQUFBzAChlxodHRwczovL2ZpZG9hbGxpYW5jZS5jby5uei90cG1wa2kvTkNVLU5UQy1LRVlJRC0zNjEwNENFNDBCQkNDMUY0MEQ4NEE0QkJENTBCRTk5MDI0RDk1N0Q0LmNydDANBgkqhkiG9w0BAQsFAAOCAQEAbp-Xp9W0vyY08YUHxerc6FnFdXZ6KFuQTZ4hze60BWexCSQOee25gqOoQaQr9ufS3ImLAoV4Ifc3vKVBQvBRwMjG3pJINoWr0p2McI0F2SNclH4M0sXFYHRlmHQ2phZB6Ddd-XL8PsGyiXRI6gVacVw5ZiVEBsRrekLH-Zy25EeqS3SxaBVnEd-HZ6BGGgbflgFtyGP9fQ5YSORC-Btno_uJbmRiZm4iHiEULp9wWEWOJIOXv9tVQKsYpPg58L1_Dgc8oml1YG5a8qK3jaR77tcUgZyYy5GOk1zIsXv36f0SkmLcNTiTjrhdGVcKs2KpW5fQgm_llQ5cvhR1jlY6dFkGCDCCBgQwggPsoAMCAQICENBTpEeEh5lpTgeR7VT9oQcwDQYJKoZIhvcNAQELBQAwgb8xCzAJBgNVBAYTAlVTMQswCQYDVQQIDAJNWTESMBAGA1UEBwwJV2FrZWZpZWxkMRYwFAYDVQQKDA1GSURPIEFsbGlhbmNlMQwwCgYDVQQLDANDV0cxNjA0BgNVBAMMLUZJRE8gRmFrZSBUUE0gUm9vdCBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkgMjAxODExMC8GCSqGSIb3DQEJARYiY29uZm9ybWFuY2UtdG9vbHNAZmlkb2FsbGlhbmNlLm9yZzAeFw0xNzAyMDEwMDAwMDBaFw0zNTAxMzEyMzU5NTlaMEExPzA9BgNVBAMTNk5DVS1OVEMtS0VZSUQtMzYxMDRDRTQwQkJDQzFGNDBEODRBNEJCRDUwQkU5OTAyNEQ5NTdENDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANc-c30RpQd-_LCoiLJbXz3t_vqciOIovwjez79_DtVgi8G9Ph-tPL-lC0ueFGBMSPcKd_RDdSFe2QCYQd9e0DtiFxra-uWGa0olI1hHI7bK2GzNAZSTKEbwgqpf8vXMQ-7SPajg6PfxSOLH_Nj2yd6tkNkUSdlGtWfY8XGB3n-q--nt3UHdUQWEtgUoTe5abBXsG7MQSuTNoad3v6vk-tLd0W44ivM6pbFqFUHchx8mGLApCpjlVXrfROaCoc9E91hG9B-WNvekJ0dM6kJ658Hy7yscQ6JdqIEolYojCtWaWNmwcfv--OE1Ax_4Ub24gl3hpB9EOcBCzpb4UFmLYUECAwEAAaOCAXcwggFzMAsGA1UdDwQEAwIBhjAWBgNVHSAEDzANMAsGCSsGAQQBgjcVHzAbBgNVHSUEFDASBgkrBgEEAYI3FSQGBWeBBQgDMBIGA1UdEwEB_wQIMAYBAf8CAQAwHwYDVR0OBBgEFsIUUX8iyOZTpXzVTN0wWH4w_2c2jF0wHwYDVR0jBBgwFqAUXH82LZCtWry6jnXa3jqg7cFOAoswaAYDVR0fBGEwXzBdoFugWYZXaHR0cHM6Ly9maWRvYWxsaWFuY2UuY28ubnovdHBtcGtpL2NybC9GSURPIEZha2UgVFBNIFJvb3QgQ2VydGlmaWNhdGUgQXV0aG9yaXR5IDIwMTguY3JsMG8GCCsGAQUFBwEBBGMwYTBfBggrBgEFBQcwAoZTaHR0cHM6Ly9maWRvYWxsaWFuY2UuY28ubnovdHBtcGtpL0ZJRE8gRmFrZSBUUE0gUm9vdCBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkgMjAxOC5jcnQwDQYJKoZIhvcNAQELBQADggIBAG138t55DF9nPJbvbPQZOypmyTPpNne0A5fh69P1fHZ5qdE2PDz3cf5Tl-8OPI4xQniEFNPcXMb7KlhMM6zCl4GkZtNN4MxygdFjQ1gTZOBDpt7Dwziij0MakmwyC0RYTNtbSyVhHUevgw9rnu13EzqxPyL5JD-UqADh2Y51MS0qy7IOgegLQv-eJzSNUgHxFJreUzz4PU6yzSsTyyYDW-H4ZjAQKienVp8ewZf8oHGWHGQFGa5E9m1P8vxCMZ7pIzeQweCVYrs3q7unu4nzBAIXLPI092kYFUgyz3lIaSB3XEiPBokpupX6Zmgrfphb-XX3tbenH5hkxfumueA5RMHTMu5TVjhJXiV0yM3q5W5xrQHdJlF5nOdJDEE-Kb7nm6xaT1DDpafqBc5vEDMkJmBA4AXHUY7JPGqEEzEenT7k6Wn5IQLZg4qc8Irnj__yM7xUhJWJam47KVbLA4WFu-IKvJrkP5GSglZ9qASOCxBHaOL2UcTAg50uvhUSwur2KSak2vlENdmAijwdAL4LLQWrkFd-9NBwcNwTdfK4ekEHP1l4BwJtkNwW6etUgeA5rkW2JLocXoBq5v7GSk4_CBoKhyiahQGQQ9SZFGeBJhzzkK9yN-yKskcVjjjInSHPl-ZpeOK3sI08sEyTH0gxlTtRoX0MKDsMAHEVToe5o1u9Z3B1YkFyZWFZATYAAQALAAYEcgAgnf_L82w4OuaZ-5ho3G3LidcVOIS-KAOSLBJBWL-tIq4AEAAQCAAAAAAAAQDPtSggWlsjcFiQO61-hUF8i-3FPcyvuARcy3p1seZ-_B4ClhNh5U-T0v0flMU5p6nsNDWj4f6-soe-2vVJMTm2d26uKYD2zwdrkrYYXRu5IFqUXqF-kY99v8RcrAF7DQKDo-E4XhiMz6uECvnjEloGfTYZrVuQ1mdjQ8Qki7U-9SQHMW_IsaI8ZKHtupXNhM5YPQyFbDHHXSE_iyPGh2mY4SR466ouesIuG0NccCUk5UDIvS__OUmNaX7aBrKTlnkMFjkCA1ZDFC99ZQoLFCJQHqnOU7m8zSvTJpUyG2feWgAL2Gl05V3I_lb_v5yELXcihFoA33QIOSpDmKqKV3SXaGNlcnRJbmZvWK3_VENHgBcAIgALEeaO1E21Ny4UKW4vhKzHg5h1GIGSHjD8IqBvi3PHlFMAIBo8rAwJFDGsmQjauX_FCBQenvBa2ApBcR_gOx2qW2QAAAAAAUdwF0hVaXtLxoVgpQFzfvmNNFZV-wAiAAsXPoJSq0uhvU6VLf0uIelHBNFHEanasKAoTp-lQ2dRGAAiAAuO1HPzTRRabZhwPvHQh0b1MnLIG8EVGNfpshASWSfjQWhhdXRoRGF0YVkBZ0mWDeWIDoxodDQXD2R2YFuP5K65ooYyx5lc87qDHZdjQQAAAEOn1tk6ig0R6JqUps9xBy9zACCH1cyGRV483U-ur0qz9V_AixVm-36OZJFMSd69Nz4oH6QBAwM5AQAgWQEAz7UoIFpbI3BYkDutfoVBfIvtxT3Mr7gEXMt6dbHmfvweApYTYeVPk9L9H5TFOaep7DQ1o-H-vrKHvtr1STE5tndurimA9s8Ha5K2GF0buSBalF6hfpGPfb_EXKwBew0Cg6PhOF4YjM-rhAr54xJaBn02Ga1bkNZnY0PEJIu1PvUkBzFvyLGiPGSh7bqVzYTOWD0MhWwxx10hP4sjxodpmOEkeOuqLnrCLhtDXHAlJOVAyL0v_zlJjWl-2gayk5Z5DBY5AgNWQxQvfWUKCxQiUB6pzlO5vM0r0yaVMhtn3loAC9hpdOVdyP5W_7-chC13IoRaAN90CDkqQ5iqild0lyFDAQAB"),
	},
}

func TestTPMAttestationVerificationFailAttStatement(t *testing.T) {
	tests := []struct {
		name    string
		att     types.AttestationObject
		wantErr string
	}{
		{
			"TPM Negative Test AttStatement Missing Ver",
			types.AttestationObject{},
			"Error retreiving ver value",
		},
		{
			"TPM Negative Test AttStatement Ver not 2.0",
			types.AttestationObject{AttStatement: map[string]interface{}{"ver": "foo.bar"}},
			"WebAuthn only supports TPM 2.0 currently",
		},
		{
			"TPM Negative Test AttStatement Alg not present",
			types.AttestationObject{AttStatement: map[string]interface{}{"ver": "2.0"}},
			"Error retreiving alg value",
		},
		{
			"TPM Negative Test AttStatement x5c not present",
			types.AttestationObject{AttStatement: map[string]interface{}{"ver": "2.0", "alg": int64(0)}},
			errors.ErrNotImplemented.Message(),
		},
		{
			"TPM Negative Test AttStatement ecdaaKeyId present",
			types.AttestationObject{AttStatement: map[string]interface{}{"ver": "2.0", "alg": int64(0), "x5c": []interface{}{}, "ecdaaKeyId": []byte{}}},
			errors.ErrNotImplemented.Message(),
		},
		{
			"TPM Negative Test AttStatement sig not present",
			types.AttestationObject{AttStatement: map[string]interface{}{"ver": "2.0", "alg": int64(0), "x5c": []interface{}{}}},
			"Error retreiving sig value",
		},
		{
			"TPM Negative Test AttStatement certInfo not present",
			types.AttestationObject{AttStatement: map[string]interface{}{"ver": "2.0", "alg": int64(0), "x5c": []interface{}{}, "sig": []byte{}}},
			"Error retreiving certInfo value",
		},
		{
			"TPM Negative Test AttStatement pubArea not present",
			types.AttestationObject{AttStatement: map[string]interface{}{"ver": "2.0", "alg": int64(0), "x5c": []interface{}{}, "sig": []byte{}, "certInfo": []byte{}}},
			"Error retreiving pubArea value",
		},
		{
			"TPM Negative Test pubArea not TPMT_PUBLIC",
			types.AttestationObject{AttStatement: defaultAttStatement},
			"Unable to decode TPMT_PUBLIC in attestation statement",
		},
	}
	for _, tt := range tests {
		attestationKey := provider.ID()
		_, _, _, err := provider.Handler(tt.att, nil)
		if tt.wantErr != "" {
			assert.Contains(t, err.Error(), tt.wantErr)
		} else {
			assert.Equal(t, "tpm", attestationKey)
		}
	}
}

var defaultAttStatement = map[string]interface{}{"ver": "2.0", "alg": int64(0), "x5c": []interface{}{}, "sig": []byte{}, "certInfo": []byte{}, "pubArea": []byte{}}

func TestTPMAttestationVerificationFailPubArea(t *testing.T) {
	tests := []struct {
		name    string
		public  googletpm.Public
		wantErr string
	}{}
	for _, tt := range tests {
		attStmt := make(map[string]interface{}, len(defaultAttStatement))
		for id, v := range defaultAttStatement {
			attStmt[id] = v
		}

		//attStmt["pubArea"], _ = tt.public.Encode()
		att := types.AttestationObject{
			AttStatement: attStmt,
		}
		_, attestationKey, _, err := provider.Handler(att, nil)
		if tt.wantErr != "" {
			assert.Contains(t, err.Error(), tt.wantErr)
		} else {
			assert.Equal(t, "tpm", attestationKey)
		}
	}
}

func TestTPMAttestationVerificationFailCertInfo(t *testing.T) {
	tests := []struct {
		name           string
		att            types.AttestationObject
		clientDataHash [32]byte
		wantErr        string
	}{}
	for _, tt := range tests {
		_, attestationKey, _, err := provider.Handler(tt.att, tt.clientDataHash[:])
		if tt.wantErr != "" {
			assert.Contains(t, err.Error(), tt.wantErr)
		} else {
			assert.Equal(t, "tpm", attestationKey)
		}
	}
}

func TestTPMAttestationVerificationFailX5c(t *testing.T) {
	tests := []struct {
		name           string
		att            types.AttestationObject
		clientDataHash [32]byte
		wantErr        string
	}{}
	for _, tt := range tests {
		_, attestationKey, _, err := provider.Handler(tt.att, tt.clientDataHash[:])
		if tt.wantErr != "" {
			assert.Contains(t, err.Error(), tt.wantErr)
		} else {
			assert.Equal(t, "tpm", attestationKey)
		}
	}
}
