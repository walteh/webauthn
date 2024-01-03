package awsrds


// The type of storage.
//
// Example:
//   // Set open cursors with parameter group
//   parameterGroup := rds.NewParameterGroup(this, jsii.String("ParameterGroup"), &ParameterGroupProps{
//   	Engine: rds.DatabaseInstanceEngine_OracleSe2(&OracleSe2InstanceEngineProps{
//   		Version: rds.OracleEngineVersion_VER_19_0_0_0_2020_04_R1(),
//   	}),
//   	Parameters: map[string]*string{
//   		"open_cursors": jsii.String("2500"),
//   	},
//   })
//
//   optionGroup := rds.NewOptionGroup(this, jsii.String("OptionGroup"), &OptionGroupProps{
//   	Engine: rds.DatabaseInstanceEngine_*OracleSe2(&OracleSe2InstanceEngineProps{
//   		Version: rds.OracleEngineVersion_VER_19_0_0_0_2020_04_R1(),
//   	}),
//   	Configurations: []optionConfiguration{
//   		&optionConfiguration{
//   			Name: jsii.String("LOCATOR"),
//   		},
//   		&optionConfiguration{
//   			Name: jsii.String("OEM"),
//   			Port: jsii.Number(1158),
//   			Vpc: *Vpc,
//   		},
//   	},
//   })
//
//   // Allow connections to OEM
//   optionGroup.OptionConnections.oEM.Connections.AllowDefaultPortFromAnyIpv4()
//
//   // Database instance with production values
//   instance := rds.NewDatabaseInstance(this, jsii.String("Instance"), &DatabaseInstanceProps{
//   	Engine: rds.DatabaseInstanceEngine_*OracleSe2(&OracleSe2InstanceEngineProps{
//   		Version: rds.OracleEngineVersion_VER_19_0_0_0_2020_04_R1(),
//   	}),
//   	LicenseModel: rds.LicenseModel_BRING_YOUR_OWN_LICENSE,
//   	InstanceType: ec2.InstanceType_Of(ec2.InstanceClass_BURSTABLE3, ec2.InstanceSize_MEDIUM),
//   	MultiAz: jsii.Boolean(true),
//   	StorageType: rds.StorageType_IO1,
//   	Credentials: rds.Credentials_FromUsername(jsii.String("syscdk")),
//   	Vpc: Vpc,
//   	DatabaseName: jsii.String("ORCL"),
//   	StorageEncrypted: jsii.Boolean(true),
//   	BackupRetention: cdk.Duration_Days(jsii.Number(7)),
//   	MonitoringInterval: cdk.Duration_Seconds(jsii.Number(60)),
//   	EnablePerformanceInsights: jsii.Boolean(true),
//   	CloudwatchLogsExports: []*string{
//   		jsii.String("trace"),
//   		jsii.String("audit"),
//   		jsii.String("alert"),
//   		jsii.String("listener"),
//   	},
//   	CloudwatchLogsRetention: logs.RetentionDays_ONE_MONTH,
//   	AutoMinorVersionUpgrade: jsii.Boolean(true),
//   	 // required to be true if LOCATOR is used in the option group
//   	OptionGroup: OptionGroup,
//   	ParameterGroup: ParameterGroup,
//   	RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
//   })
//
//   // Allow connections on default port from any IPV4
//   instance.connections.AllowDefaultPortFromAnyIpv4()
//
//   // Rotate the master user password every 30 days
//   instance.addRotationSingleUser()
//
//   // Add alarm for high CPU
//   // Add alarm for high CPU
//   cloudwatch.NewAlarm(this, jsii.String("HighCPU"), &AlarmProps{
//   	Metric: instance.metricCPUUtilization(),
//   	Threshold: jsii.Number(90),
//   	EvaluationPeriods: jsii.Number(1),
//   })
//
//   // Trigger Lambda function on instance availability events
//   fn := lambda.NewFunction(this, jsii.String("Function"), &FunctionProps{
//   	Code: lambda.Code_FromInline(jsii.String("exports.handler = (event) => console.log(event);")),
//   	Handler: jsii.String("index.handler"),
//   	Runtime: lambda.Runtime_NODEJS_18_X(),
//   })
//
//   availabilityRule := instance.OnEvent(jsii.String("Availability"), &OnEventOptions{
//   	Target: targets.NewLambdaFunction(fn),
//   })
//   availabilityRule.AddEventPattern(&EventPattern{
//   	Detail: map[string]interface{}{
//   		"EventCategories": []interface{}{
//   			jsii.String("availability"),
//   		},
//   	},
//   })
//
// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html
//
type StorageType string

const (
	// Standard.
	//
	// Amazon RDS supports magnetic storage for backward compatibility. It is recommended to use
	// General Purpose SSD or Provisioned IOPS SSD for any new storage needs.
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html#CHAP_Storage.Magnetic
	//
	StorageType_STANDARD StorageType = "STANDARD"
	// General purpose SSD (gp2).
	//
	// Baseline performance determined by volume size.
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html#Concepts.Storage.GeneralSSD
	//
	StorageType_GP2 StorageType = "GP2"
	// General purpose SSD (gp3).
	//
	// Performance scales independently from storage.
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html#Concepts.Storage.GeneralSSD
	//
	StorageType_GP3 StorageType = "GP3"
	// Provisioned IOPS (SSD).
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html#USER_PIOPS
	//
	StorageType_IO1 StorageType = "IO1"
)

