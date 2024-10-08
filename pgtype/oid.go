package pgtype

// Oid is a Postgres Object ID.
type Oid uint32

// String 实现 Stringer 接口
func (o Oid) String() string {
	return TypeName[o]
}

const (
	T_bool             Oid = 16
	T_bytea            Oid = 17
	T_char             Oid = 18
	T_name             Oid = 19
	T_int8             Oid = 20
	T_int2             Oid = 21
	T_int2vector       Oid = 22
	T_int4             Oid = 23
	T_regproc          Oid = 24
	T_text             Oid = 25
	T_oid              Oid = 26
	T_tid              Oid = 27
	T_xid              Oid = 28
	T_cid              Oid = 29
	T_oidvector        Oid = 30
	T_pg_ddl_command   Oid = 32
	T_pg_type          Oid = 71
	T_pg_attribute     Oid = 75
	T_pg_proc          Oid = 81
	T_pg_class         Oid = 83
	T_json             Oid = 114
	T_xml              Oid = 142
	T__xml             Oid = 143
	T_pg_node_tree     Oid = 194
	T__json            Oid = 199
	T_smgr             Oid = 210
	T_index_am_handler Oid = 325
	T_point            Oid = 600
	T_lseg             Oid = 601
	T_path             Oid = 602
	T_box              Oid = 603
	T_polygon          Oid = 604
	T_line             Oid = 628
	T__line            Oid = 629
	T_cidr             Oid = 650
	T__cidr            Oid = 651
	T_float4           Oid = 700
	T_float8           Oid = 701
	T_abstime          Oid = 702
	T_reltime          Oid = 703
	T_tinterval        Oid = 704
	T_unknown          Oid = 705
	T_circle           Oid = 718
	T__circle          Oid = 719
	T_money            Oid = 790
	T__money           Oid = 791
	T_macaddr          Oid = 829
	T_inet             Oid = 869
	T__bool            Oid = 1000
	T__bytea           Oid = 1001
	T__char            Oid = 1002
	T__name            Oid = 1003
	T__int2            Oid = 1005
	T__int2vector      Oid = 1006
	T__int4            Oid = 1007
	T__regproc         Oid = 1008
	T__text            Oid = 1009
	T__tid             Oid = 1010
	T__xid             Oid = 1011
	T__cid             Oid = 1012
	T__oidvector       Oid = 1013
	T__bpchar          Oid = 1014
	T__varchar         Oid = 1015
	T__int8            Oid = 1016
	T__point           Oid = 1017
	T__lseg            Oid = 1018
	T__path            Oid = 1019
	T__box             Oid = 1020
	T__float4          Oid = 1021
	T__float8          Oid = 1022
	T__abstime         Oid = 1023
	T__reltime         Oid = 1024
	T__tinterval       Oid = 1025
	T__polygon         Oid = 1027
	T__oid             Oid = 1028
	T_aclitem          Oid = 1033
	T__aclitem         Oid = 1034
	T__macaddr         Oid = 1040
	T__inet            Oid = 1041
	T_bpchar           Oid = 1042
	T_varchar          Oid = 1043
	T_date             Oid = 1082
	T_time             Oid = 1083
	T_timestamp        Oid = 1114
	T__timestamp       Oid = 1115
	T__date            Oid = 1182
	T__time            Oid = 1183
	T_timestamptz      Oid = 1184
	T__timestamptz     Oid = 1185
	T_interval         Oid = 1186
	T__interval        Oid = 1187
	T__numeric         Oid = 1231
	T_pg_database      Oid = 1248
	T__cstring         Oid = 1263
	T_timetz           Oid = 1266
	T__timetz          Oid = 1270
	T_bit              Oid = 1560
	T__bit             Oid = 1561
	T_varbit           Oid = 1562
	T__varbit          Oid = 1563
	T_numeric          Oid = 1700
	T_refcursor        Oid = 1790
	T__refcursor       Oid = 2201
	T_regprocedure     Oid = 2202
	T_regoper          Oid = 2203
	T_regoperator      Oid = 2204
	T_regclass         Oid = 2205
	T_regtype          Oid = 2206
	T__regprocedure    Oid = 2207
	T__regoper         Oid = 2208
	T__regoperator     Oid = 2209
	T__regclass        Oid = 2210
	T__regtype         Oid = 2211
	T_record           Oid = 2249
	T_cstring          Oid = 2275
	T_any              Oid = 2276
	T_anyarray         Oid = 2277
	T_void             Oid = 2278
	T_trigger          Oid = 2279
	T_language_handler Oid = 2280
	T_internal         Oid = 2281
	T_opaque           Oid = 2282
	T_anyelement       Oid = 2283
	T__record          Oid = 2287
	T_anynonarray      Oid = 2776
	T_pg_authid        Oid = 2842
	T_pg_auth_members  Oid = 2843
	T__txid_snapshot   Oid = 2949
	T_uuid             Oid = 2950
	T__uuid            Oid = 2951
	T_txid_snapshot    Oid = 2970
	T_fdw_handler      Oid = 3115
	T_pg_lsn           Oid = 3220
	T__pg_lsn          Oid = 3221
	T_tsm_handler      Oid = 3310
	T_anyenum          Oid = 3500
	T_tsvector         Oid = 3614
	T_tsquery          Oid = 3615
	T_gtsvector        Oid = 3642
	T__tsvector        Oid = 3643
	T__gtsvector       Oid = 3644
	T__tsquery         Oid = 3645
	T_regconfig        Oid = 3734
	T__regconfig       Oid = 3735
	T_regdictionary    Oid = 3769
	T__regdictionary   Oid = 3770
	T_jsonb            Oid = 3802
	T__jsonb           Oid = 3807
	T_anyrange         Oid = 3831
	T_event_trigger    Oid = 3838
	T_int4range        Oid = 3904
	T__int4range       Oid = 3905
	T_numrange         Oid = 3906
	T__numrange        Oid = 3907
	T_tsrange          Oid = 3908
	T__tsrange         Oid = 3909
	T_tstzrange        Oid = 3910
	T__tstzrange       Oid = 3911
	T_daterange        Oid = 3912
	T__daterange       Oid = 3913
	T_int8range        Oid = 3926
	T__int8range       Oid = 3927
	T_pg_shseclabel    Oid = 4066
	T_regnamespace     Oid = 4089
	T__regnamespace    Oid = 4090
	T_regrole          Oid = 4096
	T__regrole         Oid = 4097
)

var TypeName = map[Oid]string{
	T_bool:             "BOOL",
	T_bytea:            "BYTEA",
	T_char:             "CHAR",
	T_name:             "NAME",
	T_int8:             "INT8",
	T_int2:             "INT2",
	T_int2vector:       "INT2VECTOR",
	T_int4:             "INT4",
	T_regproc:          "REGPROC",
	T_text:             "TEXT",
	T_oid:              "OID",
	T_tid:              "TID",
	T_xid:              "XID",
	T_cid:              "CID",
	T_oidvector:        "OIDVECTOR",
	T_pg_ddl_command:   "PG_DDL_COMMAND",
	T_pg_type:          "PG_TYPE",
	T_pg_attribute:     "PG_ATTRIBUTE",
	T_pg_proc:          "PG_PROC",
	T_pg_class:         "PG_CLASS",
	T_json:             "JSON",
	T_xml:              "XML",
	T__xml:             "_XML",
	T_pg_node_tree:     "PG_NODE_TREE",
	T__json:            "_JSON",
	T_smgr:             "SMGR",
	T_index_am_handler: "INDEX_AM_HANDLER",
	T_point:            "POINT",
	T_lseg:             "LSEG",
	T_path:             "PATH",
	T_box:              "BOX",
	T_polygon:          "POLYGON",
	T_line:             "LINE",
	T__line:            "_LINE",
	T_cidr:             "CIDR",
	T__cidr:            "_CIDR",
	T_float4:           "FLOAT4",
	T_float8:           "FLOAT8",
	T_abstime:          "ABSTIME",
	T_reltime:          "RELTIME",
	T_tinterval:        "TINTERVAL",
	T_unknown:          "UNKNOWN",
	T_circle:           "CIRCLE",
	T__circle:          "_CIRCLE",
	T_money:            "MONEY",
	T__money:           "_MONEY",
	T_macaddr:          "MACADDR",
	T_inet:             "INET",
	T__bool:            "_BOOL",
	T__bytea:           "_BYTEA",
	T__char:            "_CHAR",
	T__name:            "_NAME",
	T__int2:            "_INT2",
	T__int2vector:      "_INT2VECTOR",
	T__int4:            "_INT4",
	T__regproc:         "_REGPROC",
	T__text:            "_TEXT",
	T__tid:             "_TID",
	T__xid:             "_XID",
	T__cid:             "_CID",
	T__oidvector:       "_OIDVECTOR",
	T__bpchar:          "_BPCHAR",
	T__varchar:         "_VARCHAR",
	T__int8:            "_INT8",
	T__point:           "_POINT",
	T__lseg:            "_LSEG",
	T__path:            "_PATH",
	T__box:             "_BOX",
	T__float4:          "_FLOAT4",
	T__float8:          "_FLOAT8",
	T__abstime:         "_ABSTIME",
	T__reltime:         "_RELTIME",
	T__tinterval:       "_TINTERVAL",
	T__polygon:         "_POLYGON",
	T__oid:             "_OID",
	T_aclitem:          "ACLITEM",
	T__aclitem:         "_ACLITEM",
	T__macaddr:         "_MACADDR",
	T__inet:            "_INET",
	T_bpchar:           "BPCHAR",
	T_varchar:          "VARCHAR",
	T_date:             "DATE",
	T_time:             "TIME",
	T_timestamp:        "TIMESTAMP",
	T__timestamp:       "_TIMESTAMP",
	T__date:            "_DATE",
	T__time:            "_TIME",
	T_timestamptz:      "TIMESTAMPTZ",
	T__timestamptz:     "_TIMESTAMPTZ",
	T_interval:         "INTERVAL",
	T__interval:        "_INTERVAL",
	T__numeric:         "_NUMERIC",
	T_pg_database:      "PG_DATABASE",
	T__cstring:         "_CSTRING",
	T_timetz:           "TIMETZ",
	T__timetz:          "_TIMETZ",
	T_bit:              "BIT",
	T__bit:             "_BIT",
	T_varbit:           "VARBIT",
	T__varbit:          "_VARBIT",
	T_numeric:          "NUMERIC",
	T_refcursor:        "REFCURSOR",
	T__refcursor:       "_REFCURSOR",
	T_regprocedure:     "REGPROCEDURE",
	T_regoper:          "REGOPER",
	T_regoperator:      "REGOPERATOR",
	T_regclass:         "REGCLASS",
	T_regtype:          "REGTYPE",
	T__regprocedure:    "_REGPROCEDURE",
	T__regoper:         "_REGOPER",
	T__regoperator:     "_REGOPERATOR",
	T__regclass:        "_REGCLASS",
	T__regtype:         "_REGTYPE",
	T_record:           "RECORD",
	T_cstring:          "CSTRING",
	T_any:              "ANY",
	T_anyarray:         "ANYARRAY",
	T_void:             "VOID",
	T_trigger:          "TRIGGER",
	T_language_handler: "LANGUAGE_HANDLER",
	T_internal:         "INTERNAL",
	T_opaque:           "OPAQUE",
	T_anyelement:       "ANYELEMENT",
	T__record:          "_RECORD",
	T_anynonarray:      "ANYNONARRAY",
	T_pg_authid:        "PG_AUTHID",
	T_pg_auth_members:  "PG_AUTH_MEMBERS",
	T__txid_snapshot:   "_TXID_SNAPSHOT",
	T_uuid:             "UUID",
	T__uuid:            "_UUID",
	T_txid_snapshot:    "TXID_SNAPSHOT",
	T_fdw_handler:      "FDW_HANDLER",
	T_pg_lsn:           "PG_LSN",
	T__pg_lsn:          "_PG_LSN",
	T_tsm_handler:      "TSM_HANDLER",
	T_anyenum:          "ANYENUM",
	T_tsvector:         "TSVECTOR",
	T_tsquery:          "TSQUERY",
	T_gtsvector:        "GTSVECTOR",
	T__tsvector:        "_TSVECTOR",
	T__gtsvector:       "_GTSVECTOR",
	T__tsquery:         "_TSQUERY",
	T_regconfig:        "REGCONFIG",
	T__regconfig:       "_REGCONFIG",
	T_regdictionary:    "REGDICTIONARY",
	T__regdictionary:   "_REGDICTIONARY",
	T_jsonb:            "JSONB",
	T__jsonb:           "_JSONB",
	T_anyrange:         "ANYRANGE",
	T_event_trigger:    "EVENT_TRIGGER",
	T_int4range:        "INT4RANGE",
	T__int4range:       "_INT4RANGE",
	T_numrange:         "NUMRANGE",
	T__numrange:        "_NUMRANGE",
	T_tsrange:          "TSRANGE",
	T__tsrange:         "_TSRANGE",
	T_tstzrange:        "TSTZRANGE",
	T__tstzrange:       "_TSTZRANGE",
	T_daterange:        "DATERANGE",
	T__daterange:       "_DATERANGE",
	T_int8range:        "INT8RANGE",
	T__int8range:       "_INT8RANGE",
	T_pg_shseclabel:    "PG_SHSECLABEL",
	T_regnamespace:     "REGNAMESPACE",
	T__regnamespace:    "_REGNAMESPACE",
	T_regrole:          "REGROLE",
	T__regrole:         "_REGROLE",
}
