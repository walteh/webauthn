package awsrds

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
)

// Construction properties for a DatabaseInstanceReadReplica.
//
// Example:
//   var vpc vpc
//
//   var sourceInstance databaseInstance
//
//   rds.NewDatabaseInstanceFromSnapshot(this, jsii.String("Instance"), &DatabaseInstanceFromSnapshotProps{
//   	SnapshotIdentifier: jsii.String("my-snapshot"),
//   	Engine: rds.DatabaseInstanceEngine_Postgres(&PostgresInstanceEngineProps{
//   		Version: rds.PostgresEngineVersion_VER_15_2(),
//   	}),
//   	// optional, defaults to m5.large
//   	InstanceType: ec2.InstanceType_Of(ec2.InstanceClass_BURSTABLE2, ec2.InstanceSize_LARGE),
//   	Vpc: Vpc,
//   })
//   rds.NewDatabaseInstanceReadReplica(this, jsii.String("ReadReplica"), &DatabaseInstanceReadReplicaProps{
//   	SourceDatabaseInstance: sourceInstance,
//   	InstanceType: ec2.InstanceType_*Of(ec2.InstanceClass_BURSTABLE2, ec2.InstanceSize_LARGE),
//   	Vpc: Vpc,
//   })
//
type DatabaseInstanceReadReplicaProps struct {
	// The VPC network where the DB subnet group should be created.
	Vpc awsec2.IVpc `field:"required" json:"vpc" yaml:"vpc"`
	// Indicates that minor engine upgrades are applied automatically to the DB instance during the maintenance window.
	// Default: true.
	//
	AutoMinorVersionUpgrade *bool `field:"optional" json:"autoMinorVersionUpgrade" yaml:"autoMinorVersionUpgrade"`
	// The name of the Availability Zone where the DB instance will be located.
	// Default: - no preference.
	//
	AvailabilityZone *string `field:"optional" json:"availabilityZone" yaml:"availabilityZone"`
	// The number of days during which automatic DB snapshots are retained.
	//
	// Set to zero to disable backups.
	// When creating a read replica, you must enable automatic backups on the source
	// database instance by setting the backup retention to a value other than zero.
	// Default: - Duration.days(1) for source instances, disabled for read replicas
	//
	BackupRetention awscdk.Duration `field:"optional" json:"backupRetention" yaml:"backupRetention"`
	// The identifier of the CA certificate for this DB instance.
	//
	// Specifying or updating this property triggers a reboot.
	//
	// For RDS DB engines:.
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/UsingWithRDS.SSL-certificate-rotation.html
	//
	// Default: - RDS will choose a certificate authority.
	//
	CaCertificate CaCertificate `field:"optional" json:"caCertificate" yaml:"caCertificate"`
	// The list of log types that need to be enabled for exporting to CloudWatch Logs.
	// Default: - no log exports.
	//
	CloudwatchLogsExports *[]*string `field:"optional" json:"cloudwatchLogsExports" yaml:"cloudwatchLogsExports"`
	// The number of days log events are kept in CloudWatch Logs.
	//
	// When updating
	// this property, unsetting it doesn't remove the log retention policy. To
	// remove the retention policy, set the value to `Infinity`.
	// Default: - logs never expire.
	//
	CloudwatchLogsRetention awslogs.RetentionDays `field:"optional" json:"cloudwatchLogsRetention" yaml:"cloudwatchLogsRetention"`
	// The IAM role for the Lambda function associated with the custom resource that sets the retention policy.
	// Default: - a new role is created.
	//
	CloudwatchLogsRetentionRole awsiam.IRole `field:"optional" json:"cloudwatchLogsRetentionRole" yaml:"cloudwatchLogsRetentionRole"`
	// Indicates whether to copy all of the user-defined tags from the DB instance to snapshots of the DB instance.
	// Default: true.
	//
	CopyTagsToSnapshot *bool `field:"optional" json:"copyTagsToSnapshot" yaml:"copyTagsToSnapshot"`
	// Indicates whether automated backups should be deleted or retained when you delete a DB instance.
	// Default: true.
	//
	DeleteAutomatedBackups *bool `field:"optional" json:"deleteAutomatedBackups" yaml:"deleteAutomatedBackups"`
	// Indicates whether the DB instance should have deletion protection enabled.
	// Default: - true if ``removalPolicy`` is RETAIN, false otherwise.
	//
	DeletionProtection *bool `field:"optional" json:"deletionProtection" yaml:"deletionProtection"`
	// The Active Directory directory ID to create the DB instance in.
	// Default: - Do not join domain.
	//
	Domain *string `field:"optional" json:"domain" yaml:"domain"`
	// The IAM role to be used when making API calls to the Directory Service.
	//
	// The role needs the AWS-managed policy
	// AmazonRDSDirectoryServiceAccess or equivalent.
	// Default: - The role will be created for you if `DatabaseInstanceNewProps#domain` is specified.
	//
	DomainRole awsiam.IRole `field:"optional" json:"domainRole" yaml:"domainRole"`
	// Whether to enable Performance Insights for the DB instance.
	// Default: - false, unless ``performanceInsightRetention`` or ``performanceInsightEncryptionKey`` is set.
	//
	EnablePerformanceInsights *bool `field:"optional" json:"enablePerformanceInsights" yaml:"enablePerformanceInsights"`
	// Whether to enable mapping of AWS Identity and Access Management (IAM) accounts to database accounts.
	// Default: false.
	//
	IamAuthentication *bool `field:"optional" json:"iamAuthentication" yaml:"iamAuthentication"`
	// A name for the DB instance.
	//
	// If you specify a name, AWS CloudFormation
	// converts it to lowercase.
	// Default: - a CloudFormation generated name.
	//
	InstanceIdentifier *string `field:"optional" json:"instanceIdentifier" yaml:"instanceIdentifier"`
	// The number of I/O operations per second (IOPS) that the database provisions.
	//
	// The value must be equal to or greater than 1000.
	// Default: - no provisioned iops if storage type is not specified. For GP3: 3,000 IOPS if allocated
	// storage is less than 400 GiB for MariaDB, MySQL, and PostgreSQL, less than 200 GiB for Oracle and
	// less than 20 GiB for SQL Server. 12,000 IOPS otherwise (except for SQL Server where the default is
	// always 3,000 IOPS).
	//
	Iops *float64 `field:"optional" json:"iops" yaml:"iops"`
	// Upper limit to which RDS can scale the storage in GiB(Gibibyte).
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_PIOPS.StorageTypes.html#USER_PIOPS.Autoscaling
	//
	// Default: - No autoscaling of RDS instance.
	//
	MaxAllocatedStorage *float64 `field:"optional" json:"maxAllocatedStorage" yaml:"maxAllocatedStorage"`
	// The interval, in seconds, between points when Amazon RDS collects enhanced monitoring metrics for the DB instance.
	// Default: - no enhanced monitoring.
	//
	MonitoringInterval awscdk.Duration `field:"optional" json:"monitoringInterval" yaml:"monitoringInterval"`
	// Role that will be used to manage DB instance monitoring.
	// Default: - A role is automatically created for you.
	//
	MonitoringRole awsiam.IRole `field:"optional" json:"monitoringRole" yaml:"monitoringRole"`
	// Specifies if the database instance is a multiple Availability Zone deployment.
	// Default: false.
	//
	MultiAz *bool `field:"optional" json:"multiAz" yaml:"multiAz"`
	// The network type of the DB instance.
	// Default: - IPV4.
	//
	NetworkType NetworkType `field:"optional" json:"networkType" yaml:"networkType"`
	// The option group to associate with the instance.
	// Default: - no option group.
	//
	OptionGroup IOptionGroup `field:"optional" json:"optionGroup" yaml:"optionGroup"`
	// The DB parameter group to associate with the instance.
	// Default: - no parameter group.
	//
	ParameterGroup IParameterGroup `field:"optional" json:"parameterGroup" yaml:"parameterGroup"`
	// The AWS KMS key for encryption of Performance Insights data.
	// Default: - default master key.
	//
	PerformanceInsightEncryptionKey awskms.IKey `field:"optional" json:"performanceInsightEncryptionKey" yaml:"performanceInsightEncryptionKey"`
	// The amount of time, in days, to retain Performance Insights data.
	// Default: 7 this is the free tier.
	//
	PerformanceInsightRetention PerformanceInsightRetention `field:"optional" json:"performanceInsightRetention" yaml:"performanceInsightRetention"`
	// The port for the instance.
	// Default: - the default port for the chosen engine.
	//
	Port *float64 `field:"optional" json:"port" yaml:"port"`
	// The daily time range during which automated backups are performed.
	//
	// Constraints:
	// - Must be in the format `hh24:mi-hh24:mi`.
	// - Must be in Universal Coordinated Time (UTC).
	// - Must not conflict with the preferred maintenance window.
	// - Must be at least 30 minutes.
	// Default: - a 30-minute window selected at random from an 8-hour block of
	// time for each AWS Region. To see the time blocks available, see
	// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_WorkingWithAutomatedBackups.html#USER_WorkingWithAutomatedBackups.BackupWindow
	//
	PreferredBackupWindow *string `field:"optional" json:"preferredBackupWindow" yaml:"preferredBackupWindow"`
	// The weekly time range (in UTC) during which system maintenance can occur.
	//
	// Format: `ddd:hh24:mi-ddd:hh24:mi`
	// Constraint: Minimum 30-minute window.
	// Default: - a 30-minute window selected at random from an 8-hour block of
	// time for each AWS Region, occurring on a random day of the week. To see
	// the time blocks available, see https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_UpgradeDBInstance.Maintenance.html#Concepts.DBMaintenance
	//
	PreferredMaintenanceWindow *string `field:"optional" json:"preferredMaintenanceWindow" yaml:"preferredMaintenanceWindow"`
	// The number of CPU cores and the number of threads per core.
	// Default: - the default number of CPU cores and threads per core for the
	// chosen instance class.
	//
	// See https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.DBInstanceClass.html#USER_ConfigureProcessor
	//
	ProcessorFeatures *ProcessorFeatures `field:"optional" json:"processorFeatures" yaml:"processorFeatures"`
	// Indicates whether the DB instance is an internet-facing instance.
	//
	// If not specified,
	// the instance's vpcSubnets will be used to determine if the instance is internet-facing
	// or not.
	// Default: - `true` if the instance's `vpcSubnets` is `subnetType: SubnetType.PUBLIC`, `false` otherwise
	//
	PubliclyAccessible *bool `field:"optional" json:"publiclyAccessible" yaml:"publiclyAccessible"`
	// The CloudFormation policy to apply when the instance is removed from the stack or replaced during an update.
	// Default: - RemovalPolicy.SNAPSHOT (remove the resource, but retain a snapshot of the data)
	//
	RemovalPolicy awscdk.RemovalPolicy `field:"optional" json:"removalPolicy" yaml:"removalPolicy"`
	// S3 buckets that you want to load data into.
	//
	// This property must not be used if `s3ExportRole` is used.
	//
	// For Microsoft SQL Server:.
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/oracle-s3-integration.html
	//
	// Default: - None.
	//
	S3ExportBuckets *[]awss3.IBucket `field:"optional" json:"s3ExportBuckets" yaml:"s3ExportBuckets"`
	// Role that will be associated with this DB instance to enable S3 export.
	//
	// This property must not be used if `s3ExportBuckets` is used.
	//
	// For Microsoft SQL Server:.
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/oracle-s3-integration.html
	//
	// Default: - New role is created if `s3ExportBuckets` is set, no role is defined otherwise.
	//
	S3ExportRole awsiam.IRole `field:"optional" json:"s3ExportRole" yaml:"s3ExportRole"`
	// S3 buckets that you want to load data from.
	//
	// This feature is only supported by the Microsoft SQL Server, Oracle, and PostgreSQL engines.
	//
	// This property must not be used if `s3ImportRole` is used.
	//
	// For Microsoft SQL Server:.
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/PostgreSQL.Procedural.Importing.html
	//
	// Default: - None.
	//
	S3ImportBuckets *[]awss3.IBucket `field:"optional" json:"s3ImportBuckets" yaml:"s3ImportBuckets"`
	// Role that will be associated with this DB instance to enable S3 import.
	//
	// This feature is only supported by the Microsoft SQL Server, Oracle, and PostgreSQL engines.
	//
	// This property must not be used if `s3ImportBuckets` is used.
	//
	// For Microsoft SQL Server:.
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/PostgreSQL.Procedural.Importing.html
	//
	// Default: - New role is created if `s3ImportBuckets` is set, no role is defined otherwise.
	//
	S3ImportRole awsiam.IRole `field:"optional" json:"s3ImportRole" yaml:"s3ImportRole"`
	// The security groups to assign to the DB instance.
	// Default: - a new security group is created.
	//
	SecurityGroups *[]awsec2.ISecurityGroup `field:"optional" json:"securityGroups" yaml:"securityGroups"`
	// The storage throughput, specified in mebibytes per second (MiBps).
	//
	// Only applicable for GP3.
	// See: https://docs.aws.amazon.com//AmazonRDS/latest/UserGuide/CHAP_Storage.html#gp3-storage
	//
	// Default: - 125 MiBps if allocated storage is less than 400 GiB for MariaDB, MySQL, and PostgreSQL,
	// less than 200 GiB for Oracle and less than 20 GiB for SQL Server. 500 MiBps otherwise (except for
	// SQL Server where the default is always 125 MiBps).
	//
	StorageThroughput *float64 `field:"optional" json:"storageThroughput" yaml:"storageThroughput"`
	// The storage type.
	//
	// Storage types supported are gp2, io1, standard.
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html#Concepts.Storage.GeneralSSD
	//
	// Default: GP2.
	//
	StorageType StorageType `field:"optional" json:"storageType" yaml:"storageType"`
	// Existing subnet group for the instance.
	// Default: - a new subnet group will be created.
	//
	SubnetGroup ISubnetGroup `field:"optional" json:"subnetGroup" yaml:"subnetGroup"`
	// The type of subnets to add to the created DB subnet group.
	// Default: - private subnets.
	//
	VpcSubnets *awsec2.SubnetSelection `field:"optional" json:"vpcSubnets" yaml:"vpcSubnets"`
	// The name of the compute and memory capacity classes.
	InstanceType awsec2.InstanceType `field:"required" json:"instanceType" yaml:"instanceType"`
	// The source database instance.
	//
	// Each DB instance can have a limited number of read replicas. For more
	// information, see https://docs.aws.amazon.com/AmazonRDS/latest/DeveloperGuide/USER_ReadRepl.html.
	SourceDatabaseInstance IDatabaseInstance `field:"required" json:"sourceDatabaseInstance" yaml:"sourceDatabaseInstance"`
	// Indicates whether the DB instance is encrypted.
	// Default: - true if storageEncryptionKey has been provided, false otherwise.
	//
	StorageEncrypted *bool `field:"optional" json:"storageEncrypted" yaml:"storageEncrypted"`
	// The KMS key that's used to encrypt the DB instance.
	// Default: - default master key if storageEncrypted is true, no key otherwise.
	//
	StorageEncryptionKey awskms.IKey `field:"optional" json:"storageEncryptionKey" yaml:"storageEncryptionKey"`
}

