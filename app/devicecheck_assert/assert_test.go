package devicecheck

// AssertionResult(challenge: "abc=", payload: "dec=", assertion: "omlzaWduYXR1cmVYRzBFAiEA8wKVrDc/oJiMCEmUUCm3cELvvVjqB47rlTrIYsH+DE4CIEIfeEH4L2mXp7Sxho9Xpvbj3yXvxI4nbFbyyQ+T2qOXcWF1dGhlbnRpY2F0b3JEYXRhWCXEH8VVrPrEUwpPpKVlwZfBLdXUTSUtMymcc2mzTPV3xkAAAAAB")

import (
	"context"
	"reflect"
	"testing"

	"github.com/walteh/webauthn/pkg/dynamo"
	"github.com/walteh/webauthn/pkg/hex"
	"github.com/walteh/webauthn/pkg/webauthn/types"

	dtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TestObject struct {
	name    string
	input   DeviceCheckAssertionInput
	body    hex.Hash
	want    DeviceCheckAssertionOutput
	puts    []dtypes.Put
	checks  []dtypes.Update
	wantErr bool
}

var testA = TestObject{
	name: "A",
	input: DeviceCheckAssertionInput{
		RawAssertionObject:   hex.HexToHash("0x7b2263726564656e7469616c5f6964223a22307866623166643061633938646361323839313736316261663937613438366337353732363930306433613934313035616661353938353735663839633437323935222c22617373657274696f6e5f6f626a656374223a223078613236393733363936373665363137343735373236353538343733303435303232313030643336313135376361323132313339633431343532646435313838396435353330666333636663356634336332613937616161636330326132646638373566393032323034663133323532376263643937626563653864383465313337636331643764313163336664303762313837353938393938313931646563383265633332356136373136313735373436383635366537343639363336313734366637323434363137343631353832356334316663353535616366616334353330613466613461353635633139376331326464356434346432353264333332393963373336396233346366353737633634303030303030303031222c2273657373696f6e5f6964223a22307833613239386361323131393463356565373932306432666663353234376436666130663333306130333863663339333365313338363032363630343330623864222c2270726f7669646572223a226170706c65222c22636c69656e745f646174615f6a736f6e223a227b5c226368616c6c656e67655c223a5c22376652396a6b7450796452706b476576715a496c735f666632564e5f6f4c534b34484e42577a724972546b5c222c5c226f726967696e5c223a5c2268747470733a2f2f6e7567672e78797a5c222c5c22747970655c223a5c22776562617574686e2e6765745c227d227d"),
		ClientDataToValidate: hex.MustBase64ToHash("aGk="),
	},
	want: DeviceCheckAssertionOutput{
		SuggestedStatusCode: 204,
		OK:                  true,
	},
	puts: []dtypes.Put{
		{
			Item: map[string]dtypes.AttributeValue{
				"challenge_id":  types.S(hex.MustBase64ToHash("7fR9jktPydRpkGevqZIls_ff2VN_oLSK4HNBWzrIrTk").Hex()),
				"session_id":    types.S("0x3a298ca21194c5ee7920d2ffc5247d6fa0f330a038cf3933e138602660430b8d"),
				"credential_id": types.S("0xfb1fd0ac98dca2891761baf97a486c75726900d3a94105afa598575f89c47295"),
				"ceremony_type": types.S("webauthn.create"),
				"created_at":    types.N("1668984054"),
				"ttl":           types.N("1668984354"),
			},
			TableName: dynamo.MockCeremonyTableName(),
		},
		{
			Item: map[string]dtypes.AttributeValue{
				"created_at":       types.N("1669414368"),
				"session_id":       types.S("0x"),
				"aaguid":           types.S("0x617070617474657374646576656c6f70"),
				"clone_warning":    types.BOOL(false),
				"public_key":       types.S("0x04bee9490389b5b36c0d4bd0676c52c46426bee73ace82f6d3c4479d6b6bec24f20ad2264f7739994e636f65f280c384aa2b70c2311741027e677db62ec80071ee"),
				"attestation_type": types.S("apple-appattest"),
				"receipt":          types.S("0x308006092a864886f70d010702a0803080020101310f300d06096086480165030402010500308006092a864886f70d010701a0802480048203e831820400301f020102020101041734343937514a534144332e78797a2e6e7567672e617070308202ea020103020101048202e0308202dc30820262a00302010202060184b0d656b9300a06082a8648ce3d040302304f3123302106035504030c1a4170706c6520417070204174746573746174696f6e204341203131133011060355040a0c0a4170706c6520496e632e3113301106035504080c0a43616c69666f726e6961301e170d3232313132343232303930375a170d3233303831373034343030375a3081913149304706035504030c4066623166643061633938646361323839313736316261663937613438366337353732363930306433613934313035616661353938353735663839633437323935311a3018060355040b0c114141412043657274696669636174696f6e31133011060355040a0c0a4170706c6520496e632e3113301106035504080c0a43616c69666f726e69613059301306072a8648ce3d020106082a8648ce3d03010703420004bee9490389b5b36c0d4bd0676c52c46426bee73ace82f6d3c4479d6b6bec24f20ad2264f7739994e636f65f280c384aa2b70c2311741027e677db62ec80071eea381e63081e3300c0603551d130101ff04023000300e0603551d0f0101ff0404030204f0307106092a864886f76364080504643062a40302010abf893003020101bf893103020100bf893203020101bf893303020101bf893419041734343937514a534144332e78797a2e6e7567672e617070a5060404736b7320bf893603020105bf893703020100bf893903020100bf893a03020100301b06092a864886f763640807040e300cbf8a7808040631362e312e31303306092a864886f76364080204263024a12204208749423ff7d8e2fbea183a2a4c2936138300528398c346332c57d80d4ddf038b300a06082a8648ce3d040302036800306502301da8920cc88f57b30c2127dbe10d7533f70fee099a7ef2a78f0ae76d7d3b12886c2edce51990511361b4194768c88ccf023100f7e2cc5cdd8f903fca505149e6ba073a2679bc8620f9645736f1181a3daf17c68e201855664d3f58977c554c26308cb53028020104020101042049e739fa222e42b6d5cb2b24522477deeead2ce1097f931c3400586cb38afe9030600201050201010458734a4b726d4f414e3974307850343858636a5a38696c78783052326932566f3555314a38516d31594d546f456d4579423079487a54385854473252776153775262675a694e435a78344230477a5063464736553165413d3d300e0201060201010406415454455354300f020107020101040773616e64626f78302002010c0201010418323032322d31312d32355432323a30393a30372e3830345a302002011502041c01010418323032332d30322d32335432323a30393a30372e3830345a000000000000a080308203ae30820354a00302010202100939b4bce90cc3a1816536372f667141300a06082a8648ce3d040302307c3130302e06035504030c274170706c65204170706c69636174696f6e20496e746567726174696f6e2043412035202d20473131263024060355040b0c1d4170706c652043657274696669636174696f6e20417574686f7269747931133011060355040a0c0a4170706c6520496e632e310b3009060355040613025553301e170d3232303431393133333330335a170d3233303531393133333330325a305a3136303406035504030c2d4170706c69636174696f6e204174746573746174696f6e2046726175642052656365697074205369676e696e6731133011060355040a0c0a4170706c6520496e632e310b30090603550406130255533059301306072a8648ce3d020106082a8648ce3d0301070342000439d4f9aa9b1cc445d65ba617acf2c084ec6f0708d59014a0e76ecf3dee3999a94c6bfb0155105555646cda8e23e026011402d07e13b9541fd8b4d657d82e9378a38201d8308201d4300c0603551d130101ff04023000301f0603551d23041830168014d917fe4b6790384b92f4dbced55780140b8f3dc9304306082b0601050507010104373035303306082b060105050730018627687474703a2f2f6f6373702e6170706c652e636f6d2f6f63737030332d616169636135673130313082011c0603551d20048201133082010f3082010b06092a864886f7636405013081fd3081c306082b060105050702023081b60c81b352656c69616e6365206f6e207468697320636572746966696361746520627920616e7920706172747920617373756d657320616363657074616e6365206f6620746865207468656e206170706c696361626c65207374616e64617264207465726d7320616e6420636f6e646974696f6e73206f66207573652c20636572746966696361746520706f6c69637920616e642063657274696669636174696f6e2070726163746963652073746174656d656e74732e303506082b060105050702011629687474703a2f2f7777772e6170706c652e636f6d2f6365727469666963617465617574686f72697479301d0603551d0e04160414fb67d30dbf73b792a6265d488d2cc11d95e273f8300e0603551d0f0101ff040403020780300f06092a864886f763640c0f04020500300a06082a8648ce3d04030203480030450221009490a0673773e72f7829367623b8dd51d7c89a09eabb00e39c6e450b05580bd0022047341a2bd13cc054a80a3aaacc3cc1457c00545318ea338d7d6dd5f60b2b872e308202f93082027fa003020102021056fb83d42bff8dc3379923b55aae6ebd300a06082a8648ce3d0403033067311b301906035504030c124170706c6520526f6f74204341202d20473331263024060355040b0c1d4170706c652043657274696669636174696f6e20417574686f7269747931133011060355040a0c0a4170706c6520496e632e310b3009060355040613025553301e170d3139303332323137353333335a170d3334303332323030303030305a307c3130302e06035504030c274170706c65204170706c69636174696f6e20496e746567726174696f6e2043412035202d20473131263024060355040b0c1d4170706c652043657274696669636174696f6e20417574686f7269747931133011060355040a0c0a4170706c6520496e632e310b30090603550406130255533059301306072a8648ce3d020106082a8648ce3d0301070342000492ce63bd7d86b1ab280a3b1ce1affb04948091acf631dfa6cb28356f444be121e557dd128d8dba827c95be49fabe33caaecd0419f12f4325faf4beb3cb837ebaa381f73081f4300f0603551d130101ff040530030101ff301f0603551d23041830168014bbb0dea15833889aa48a99debebdebafdacb24ab304606082b06010505070101043a3038303606082b06010505073001862a687474703a2f2f6f6373702e6170706c652e636f6d2f6f63737030332d6170706c65726f6f746361673330370603551d1f0430302e302ca02aa0288626687474703a2f2f63726c2e6170706c652e636f6d2f6170706c65726f6f74636167332e63726c301d0603551d0e04160414d917fe4b6790384b92f4dbced55780140b8f3dc9300e0603551d0f0101ff0404030201063010060a2a864886f7636406020304020500300a06082a8648ce3d04030303680030650231008d6fa69fa1e0e4ec5b4e738a927f3d7853988ff4da1f581ec3754afe38a84c2a831a1aaa0da6646de1b993e8d1554ced0230673b2cb4e1e8370777cbd5ec76a81a3a553b3f356ac8c5e692b0e161be804969e45f2ba96ce11102aacc61d938b7734a30820243308201c9a00302010202082dc5fc88d2c54b95300a06082a8648ce3d0403033067311b301906035504030c124170706c6520526f6f74204341202d20473331263024060355040b0c1d4170706c652043657274696669636174696f6e20417574686f7269747931133011060355040a0c0a4170706c6520496e632e310b3009060355040613025553301e170d3134303433303138313930365a170d3339303433303138313930365a3067311b301906035504030c124170706c6520526f6f74204341202d20473331263024060355040b0c1d4170706c652043657274696669636174696f6e20417574686f7269747931133011060355040a0c0a4170706c6520496e632e310b30090603550406130255533076301006072a8648ce3d020106052b810400220362000498e92f3d4072a4ed93227281131cdd1095f1c5a34e71dc1416d90ee5a6052a77647b5f4e38d3bb1c44b57ff51fb632625dc9e9845b4f304f115a00fd58580ca5f50f2c4d07471375da9797976f315ced2b9d7b203bd8b954d95e99a43a510a31a3423040301d0603551d0e04160414bbb0dea15833889aa48a99debebdebafdacb24ab300f0603551d130101ff040530030101ff300e0603551d0f0101ff040403020106300a06082a8648ce3d040303036800306502310083e9c1c4165e1a5d3418d9edeff46c0e00464bb8dfb24611c50ffde67a8ca1a66bcec203d49cf593c674b86adfaa231502306d668a10cad40dd44fcd8d433eb48a63a5336ee36dda17b7641fc85326f9886274390b175bcb51a80ce81803e7a2b22800003181fd3081fa020101308190307c3130302e06035504030c274170706c65204170706c69636174696f6e20496e746567726174696f6e2043412035202d20473131263024060355040b0c1d4170706c652043657274696669636174696f6e20417574686f7269747931133011060355040a0c0a4170706c6520496e632e310b300906035504061302555302100939b4bce90cc3a1816536372f667141300d06096086480165030402010500300a06082a8648ce3d0403020447304502205a59898ffdde0d2a1b8135006b746ea2efe9f5386c920541dbd1c9912283ef38022100a3ea3ca0e32f9025fcc878dde76d56138a39b9ff52624172a68d54672cb177cb000000000000"),
				"sign_count":       types.N("0"),
				"updated_at":       types.N("1669414368"),
				"credential_id":    types.S("0xfb1fd0ac98dca2891761baf97a486c75726900d3a94105afa598575f89c47295"),
				"credential_type":  types.S("public-key"),
			},
			TableName: dynamo.MockCredentialTableName(),
		},
	},
	checks: []dtypes.Update{
		{
			Key: map[string]dtypes.AttributeValue{
				"credential_id": types.S("0xfb1fd0ac98dca2891761baf97a486c75726900d3a94105afa598575f89c47295"),
			},
			TableName: dynamo.MockCredentialTableName(),
			ExpressionAttributeValues: map[string]dtypes.AttributeValue{
				"sign_count": types.N("1"),
			},
		},
		{
			Key: map[string]dtypes.AttributeValue{
				"challenge_id": types.S(hex.MustBase64ToHash("7fR9jktPydRpkGevqZIls_ff2VN_oLSK4HNBWzrIrTk").Hex()),
			},
			TableName: dynamo.MockCeremonyTableName(),
			ExpressionAttributeValues: map[string]dtypes.AttributeValue{
				"challenge_id": nil,
			},
		},
	},
	wantErr: false,
}

var testB = TestObject{
	name: "B",
	input: DeviceCheckAssertionInput{
		RawAssertionObject:   testA.input.RawAssertionObject,
		ClientDataToValidate: hex.MustBase64ToHash("agk="),
	},
	want: DeviceCheckAssertionOutput{
		SuggestedStatusCode: 401,
		OK:                  false,
	},
	puts: testA.puts,
	checks: []dtypes.Update{
		{
			Key: map[string]dtypes.AttributeValue{
				"challenge_id": types.S(hex.MustBase64ToHash("7fR9jktPydRpkGevqZIls_ff2VN_oLSK4HNBWzrIrTk").Hex()),
			},
			TableName:                 dynamo.MockCeremonyTableName(),
			ExpressionAttributeValues: testA.puts[0].Item,
		},
		{
			Key: map[string]dtypes.AttributeValue{
				"credential_id": types.S("0xfb1fd0ac98dca2891761baf97a486c75726900d3a94105afa598575f89c47295"),
			},
			TableName:                 dynamo.MockCredentialTableName(),
			ExpressionAttributeValues: testA.puts[1].Item,
		},
	},
	wantErr: true,
}

func TestHandler_Invoke(t *testing.T) {

	dynamo.AttachLocalDynamoServer(t)

	tests := []TestObject{
		testA, testB,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			client := dynamo.NewMockClient(t)

			dynamo.MockBatchPut(t, client, tt.puts...)

			got, err := Assert(context.Background(), client, tt.input)
			if err != nil && !tt.wantErr {
				t.Errorf("handler.Invoke() - error should always be nil - error = %v", err)
				return
			}

			if got.SuggestedStatusCode != 204 {
				if !tt.wantErr {
					t.Errorf("Handler.Invoke() error = %v, wantErr %v", err, tt.wantErr)
				}

				if got.SuggestedStatusCode != tt.want.SuggestedStatusCode {
					t.Errorf("Handler.Invoke() got = %v, want %v", got, tt.want)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.Invoke() = %v, want %v", got, tt.want)
			}

			dynamo.MockBatchCheck(t, client, tt.checks...)

		})
	}
}
