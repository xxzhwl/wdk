set
foreign_key_checks=0;
drop table if exists btexec_main_process;
create table btexec_main_process
(
    Id                     int unsigned not null primary key auto_increment comment 'ID',
    SystemKey              varchar(32) not null default '' comment '所属系统名称[对应于SystemEnName配置]',
    PreemptUuid            varchar(64) not null default '' comment '抢占标识[由LocalIP+Pid+UUID组成，当主调度进程内存中的唯一抢占标识成功注册到此字段时，说明此主调度进程具备工作权限]',
    Status                 varchar(16) not null default 'wait' comment '执行状态[wait:待运行;run:运行中]',
    HeartbeatTimestampSecs int unsigned not null default 0 comment '心跳时间戳[系统使用]',
    HeartbeatTime          datetime    not null  comment '心跳时间[人类可读]',

    unique UniqSystemKey (SystemKey)

) engine = innodb charset = utf8 comment = '主进程资源抢占表';

drop table if exists btexec_task_process;
create table btexec_task_process
(
    Id               int unsigned not null primary key auto_increment comment 'ID',
    SystemKey        varchar(32)  not null default '' comment '所属系统名称[对应于SystemEnName配置]',
    BinPath          varchar(128) not null default '' comment '程序路径[如果不提供，则认为路径与主进程相同]',
    Args             text         not null comment '参数[JSON数组形式]',
    ExecType         varchar(32)  not null default 'crontab' comment '执行类型[crontab:定时任务;daemon:常驻进程]',
    CronSpec         varchar(64)  not null default '' comment '定时任务定义[仅ExecType=crontab时有效，格式如*/3 * * * * *，支持秒级定义]',
    Memo             varchar(128) not null default '' comment '备注说明',
    CreateTime       datetime     not null  comment '创建时间',
    UpdateTime       datetime     not null  comment '更新时间',
    Updator          varchar(64)  not null default '' comment '更新人',

    LastStartTime    datetime     not null comment '最近启动时间',
    LastExecDuration double precision(9, 4) not null default 0 comment '最近执行耗时(秒)',
    Status           varchar(16)  not null default '' comment '当前状态[wait:待运行;run:运行中]',
    MainProcessInfo      varchar(32)  not null default '' comment '进程信息',

    IsEnabled        tinyint      not null default 1 comment '是否启用',

    key              IdxSystemKey(SystemKey)
) engine = innodb charset = utf8 comment = '执行服务定义表';

set
foreign_key_checks=1;