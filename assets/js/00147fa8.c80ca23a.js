"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[6625],{3905:function(e,n,t){t.d(n,{Zo:function(){return p},kt:function(){return g}});var o=t(7294);function r(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function i(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);n&&(o=o.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,o)}return t}function a(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?i(Object(t),!0).forEach((function(n){r(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):i(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function s(e,n){if(null==e)return{};var t,o,r=function(e,n){if(null==e)return{};var t,o,r={},i=Object.keys(e);for(o=0;o<i.length;o++)t=i[o],n.indexOf(t)>=0||(r[t]=e[t]);return r}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(o=0;o<i.length;o++)t=i[o],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(r[t]=e[t])}return r}var l=o.createContext({}),c=function(e){var n=o.useContext(l),t=n;return e&&(t="function"==typeof e?e(n):a(a({},n),e)),t},p=function(e){var n=c(e.components);return o.createElement(l.Provider,{value:n},e.children)},u={inlineCode:"code",wrapper:function(e){var n=e.children;return o.createElement(o.Fragment,{},n)}},d=o.forwardRef((function(e,n){var t=e.components,r=e.mdxType,i=e.originalType,l=e.parentName,p=s(e,["components","mdxType","originalType","parentName"]),d=c(t),g=r,m=d["".concat(l,".").concat(g)]||d[g]||u[g]||i;return t?o.createElement(m,a(a({ref:n},p),{},{components:t})):o.createElement(m,a({ref:n},p))}));function g(e,n){var t=arguments,r=n&&n.mdxType;if("string"==typeof e||r){var i=t.length,a=new Array(i);a[0]=d;var s={};for(var l in n)hasOwnProperty.call(n,l)&&(s[l]=n[l]);s.originalType=e,s.mdxType="string"==typeof e?e:r,a[1]=s;for(var c=2;c<i;c++)a[c]=t[c];return o.createElement.apply(null,a)}return o.createElement.apply(null,t)}d.displayName="MDXCreateElement"},4542:function(e,n,t){t.r(n),t.d(n,{frontMatter:function(){return s},contentTitle:function(){return l},metadata:function(){return c},toc:function(){return p},default:function(){return d}});var o=t(7462),r=t(3366),i=(t(7294),t(3905)),a=["components"],s={id:"gosdk",title:"Creating a Connection and Opening a Session with SDK"},l=void 0,c={unversionedId:"getting-started/pre-requisite/gosdk",id:"getting-started/pre-requisite/gosdk",isDocsHomePage:!1,title:"Creating a Connection and Opening a Session with SDK",description:"\x3c!--",source:"@site/docs/getting-started/pre-requisite/gosdk.md",sourceDirName:"getting-started/pre-requisite",slug:"/getting-started/pre-requisite/gosdk",permalink:"/orion-server/docs/getting-started/pre-requisite/gosdk",tags:[],version:"current",frontMatter:{id:"gosdk",title:"Creating a Connection and Opening a Session with SDK"},sidebar:"Documentation",previous:{title:"Orion Docker",permalink:"/orion-server/docs/getting-started/launching-one-node/docker"},next:{title:"Setting Up Command Line Utilities",permalink:"/orion-server/docs/getting-started/pre-requisite/curl"}},p=[{value:"1) Cloning the SDK Repository",id:"1-cloning-the-sdk-repository",children:[],level:2},{value:"2) Copying the Crypto Materials",id:"2-copying-the-crypto-materials",children:[],level:2},{value:"3) Creating a Connection to the Orion Cluster",id:"3-creating-a-connection-to-the-orion-cluster",children:[{value:"3.1) Source Code",id:"31-source-code",children:[],level:3},{value:"3.2) Source Code Commentary",id:"32-source-code-commentary",children:[],level:3}],level:2},{value:"4) Opening a Database Session",id:"4-opening-a-database-session",children:[{value:"4.1) Source Code",id:"41-source-code",children:[],level:3},{value:"4.2) Source Code Commentary",id:"42-source-code-commentary",children:[],level:3}],level:2}],u={toc:p};function d(e){var n=e.components,t=(0,r.Z)(e,a);return(0,i.kt)("wrapper",(0,o.Z)({},u,t,{components:n,mdxType:"MDXLayout"}),(0,i.kt)("p",null,"When we use the SDK to perform queries and transactions, the following two steps must be executed first:"),(0,i.kt)("ol",null,(0,i.kt)("li",{parentName:"ol"},"Clone the SDK"),(0,i.kt)("li",{parentName:"ol"},"Creating a connection to the Orion cluster"),(0,i.kt)("li",{parentName:"ol"},"Opening a database session with the Orion cluster")),(0,i.kt)("p",null,"Let's look at these three steps."),(0,i.kt)("div",{className:"admonition admonition-info alert alert--info"},(0,i.kt)("div",{parentName:"div",className:"admonition-heading"},(0,i.kt)("h5",{parentName:"div"},(0,i.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,i.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"14",height:"16",viewBox:"0 0 14 16"},(0,i.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"}))),"info")),(0,i.kt)("div",{parentName:"div",className:"admonition-content"},(0,i.kt)("p",{parentName:"div"}," We have an example of creating a connection and opening a session at ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/hyperledger-labs/orion-sdk-go/tree/main/examples/api"},"orion-sdk-go/examples/api/"),"."))),(0,i.kt)("h2",{id:"1-cloning-the-sdk-repository"},"1) Cloning the SDK Repository"),(0,i.kt)("p",null,"To write queries and transactions using the SDK, first, execute the following steps:"),(0,i.kt)("ol",null,(0,i.kt)("li",{parentName:"ol"},"Create the required directory using the command ",(0,i.kt)("inlineCode",{parentName:"li"},"mkdir -p github.com/hyperledger-labs")),(0,i.kt)("li",{parentName:"ol"},"Change the current working directory to the above created directory by issing the command ",(0,i.kt)("inlineCode",{parentName:"li"},"cd github.com/hyperledger-labs")),(0,i.kt)("li",{parentName:"ol"},"Clone the go SDK repository with ",(0,i.kt)("inlineCode",{parentName:"li"},"git clone https://github.com/hyperledger-labs/orion-sdk-go"))),(0,i.kt)("p",null,"Then, we can use APIs provided by the SDK."),(0,i.kt)("h2",{id:"2-copying-the-crypto-materials"},"2) Copying the Crypto Materials"),(0,i.kt)("p",null,"We need root CA certificates and user certificates to submit queries and transactions using the SDK."),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"While creating a connection, we need to provide ",(0,i.kt)("em",{parentName:"li"},"RootCAs")," configuration."),(0,i.kt)("li",{parentName:"ul"},"While opening a session, we need to provide the ",(0,i.kt)("em",{parentName:"li"},"user's certificate")," and ",(0,i.kt)("em",{parentName:"li"},"private key"),".")),(0,i.kt)("p",null,"For all examples shown in this documentation, we use the crypto materials availabe at the ",(0,i.kt)("inlineCode",{parentName:"p"},"deployment/crypto")," folder in\nthe ",(0,i.kt)("inlineCode",{parentName:"p"},"orion-server")," repository."),(0,i.kt)("p",null,"Hence, copy the ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/hyperledger-labs/orion-server/tree/main/deployment/crypto"},(0,i.kt)("inlineCode",{parentName:"a"},"github.com/hyperledger-labs/orion-server/deployment/crypto")),"\nto the location where you write/use example code provided in this documentation."),(0,i.kt)("h2",{id:"3-creating-a-connection-to-the-orion-cluster"},"3) Creating a Connection to the Orion Cluster"),(0,i.kt)("h3",{id:"31-source-code"},"3.1) Source Code"),(0,i.kt)("p",null,"The following function creates a connection to our ",(0,i.kt)("a",{parentName:"p",href:"./../launching-one-node/binary"},"single node Orion cluster")," deployed using the sample configuration."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go",metastring:'title="create-connection.go"',title:'"create-connection.go"'},'package main\n\nimport (\n    "github.com/hyperledger-labs/orion-sdk-go/pkg/bcdb"\n    "github.com/hyperledger-labs/orion-sdk-go/pkg/config"\n    "github.com/hyperledger-labs/orion-server/pkg/logger"\n)\n\nfunc createConnection() (bcdb.BCDB, error) {\n    logger, err := logger.New(\n        &logger.Config{\n            Level:         "debug",\n            OutputPath:    []string{"stdout"},\n            ErrOutputPath: []string{"stderr"},\n            Encoding:      "console",\n            Name:          "bcdb-client",\n        },\n    )\n    if err != nil {\n        return nil, err\n    }\n\n    conConf := &config.ConnectionConfig{\n        ReplicaSet: []*config.Replica{\n            {\n                ID:       "bdb-node-1",\n                Endpoint: "http://127.0.0.1:6001",\n            },\n        },\n        RootCAs: []string{\n            "./crypto/CA/CA.pem",\n        },\n        Logger: logger,\n    }\n\n    db, err := bcdb.Create(conConf)\n    if err != nil {\n        return nil, err\n    }\n\n    return db, nil\n}\n')),(0,i.kt)("h3",{id:"32-source-code-commentary"},"3.2) Source Code Commentary"),(0,i.kt)("p",null,"The ",(0,i.kt)("inlineCode",{parentName:"p"},"bcdb.Create()")," method in the ",(0,i.kt)("inlineCode",{parentName:"p"},"bcdb")," package at the SDK prepares a connection context to the Orion cluster\nand loads the certificate of root certificate authorities."),(0,i.kt)("p",null,"The signature of the ",(0,i.kt)("inlineCode",{parentName:"p"},"Create()")," function is shown below:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"func Create(config *config.ConnectionConfig) (BCDB, error)\n")),(0,i.kt)("p",null,"The parameter ",(0,i.kt)("inlineCode",{parentName:"p"},"config.ConnectionConfig")," holds "),(0,i.kt)("ol",null,(0,i.kt)("li",{parentName:"ol"},"the ",(0,i.kt)("inlineCode",{parentName:"li"},"ID")," and ",(0,i.kt)("inlineCode",{parentName:"li"},"IP address")," of each Orion node in the cluster"),(0,i.kt)("li",{parentName:"ol"},"certificate of root CAs, and"),(0,i.kt)("li",{parentName:"ol"},"a logger to log messages")),(0,i.kt)("p",null,"The structure of the ",(0,i.kt)("inlineCode",{parentName:"p"},"config.ConnectionConfig")," is shown below:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"// ConnectionConfig required configuration in order to\n// open session with BCDB instance, replica set informations\n// servers root CAs\ntype ConnectionConfig struct {\n    // List of replicas URIs client can connect to\n    ReplicaSet []*Replica\n    // Keeps path to the server's root CA\n    RootCAs []string\n    // Logger instance, if nil an internal logger is created\n    Logger *logger.SugarLogger\n}\n\n// Replica\ntype Replica struct {\n    // ID replica's ID\n    ID string\n    // Endpoint the URI of the replica to connect to\n    Endpoint string\n}\n")),(0,i.kt)("p",null,"In our ",(0,i.kt)("a",{parentName:"p",href:"./../launching-one-node/binary"},"simple deployment"),", we have only one node in the cluster. Hence, we have one ",(0,i.kt)("inlineCode",{parentName:"p"},"Replica")," with the\n",(0,i.kt)("inlineCode",{parentName:"p"},"ID")," as ",(0,i.kt)("inlineCode",{parentName:"p"},"bdb-node-1")," and ",(0,i.kt)("inlineCode",{parentName:"p"},"Endpoint")," as ",(0,i.kt)("inlineCode",{parentName:"p"},"http://127.0.0.1:6001"),". Further, we have only one root certificate authority and hence, the\n",(0,i.kt)("inlineCode",{parentName:"p"},"RootCAs")," holds the path to a single CA's certificate only."),(0,i.kt)("p",null,"The ",(0,i.kt)("inlineCode",{parentName:"p"},"Create()")," would return the ",(0,i.kt)("inlineCode",{parentName:"p"},"BCDB")," implementation that allows the user to create database sessions with the Orion cluster."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"type BCDB interface {\n    // Session instantiates session to the database\n    Session(config *config.SessionConfig) (DBSession, error)\n}\n")),(0,i.kt)("h2",{id:"4-opening-a-database-session"},"4) Opening a Database Session"),(0,i.kt)("h3",{id:"41-source-code"},"4.1) Source Code"),(0,i.kt)("p",null,"Now, once we created the Orion connection and received the ",(0,i.kt)("inlineCode",{parentName:"p"},"BCDB")," object instance, we can open a database session by calling the ",(0,i.kt)("inlineCode",{parentName:"p"},"Session()")," method. The ",(0,i.kt)("inlineCode",{parentName:"p"},"Session")," object authenticates the database user against the database server.\nThe following function opens a database session for an already existing database connection."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go",metastring:'title="open-session.go"',title:'"open-session.go"'},'package main\n\nimport (\n    "time"\n\n    "github.com/hyperledger-labs/orion-sdk-go/pkg/bcdb"\n    "github.com/hyperledger-labs/orion-sdk-go/pkg/config"\n)\n\nfunc openSession(db bcdb.BCDB, userID string) (bcdb.DBSession, error) {\n    sessionConf := &config.SessionConfig{\n        UserConfig: &config.UserConfig{\n            UserID:         userID,\n            CertPath:       "./crypto/" + userID + "/" + userID + ".pem",\n            PrivateKeyPath: "./crypto/" + userID + "/" + userID + ".key",\n        },\n        TxTimeout:    20 * time.Second,\n        QueryTimeout: 10 * time.Second,\n    }\n\n    session, err := db.Session(sessionConf)\n    if err != nil {\n        return nil, err\n    }\n\n    return session, nil\n}\n')),(0,i.kt)("h3",{id:"42-source-code-commentary"},"4.2) Source Code Commentary"),(0,i.kt)("p",null,"The signature of ",(0,i.kt)("inlineCode",{parentName:"p"},"Session()")," method is shown below:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"Session(config *config.SessionConfig) (DBSession, error)\n")),(0,i.kt)("p",null,"The ",(0,i.kt)("inlineCode",{parentName:"p"},"Session()")," takes ",(0,i.kt)("inlineCode",{parentName:"p"},"config.SessionConfig")," as a parameter which holds the user configuration (user's ID and credentials) and various configuration parameters, such as transaction timeout and query timeout.\nThe structure of the ",(0,i.kt)("inlineCode",{parentName:"p"},"config.SessionConfig")," is shown below:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"// SessionConfig keeps per database session\n// configuration information\ntype SessionConfig struct {\n    UserConfig *UserConfig\n    // The transaction timeout given to the server in case of tx sync commit - `tx.Commit(true)`.\n    // SDK will wait for `TxTimeout` + some communication margin\n    // or for timeout error from server, whatever come first.\n    TxTimeout time.Duration\n    // The query timeout - SDK will wait for query result maximum `QueryTimeout` time.\n    QueryTimeout time.Duration\n}\n\n\n// UserConfig user related information\n// maintains wallet with public and private keys\ntype UserConfig struct {\n    // UserID the identity of the user\n    UserID string\n    // CertPath path to the user's certificate\n    CertPath string\n    // PrivateKeyPath path to the user's private key\n    PrivateKeyPath string\n}\n")),(0,i.kt)("p",null,"As the ",(0,i.kt)("inlineCode",{parentName:"p"},"admin")," user is submitting the transactions, we have set the ",(0,i.kt)("inlineCode",{parentName:"p"},"UserConfig")," to hold the userID of ",(0,i.kt)("inlineCode",{parentName:"p"},"admin"),", certificate, and private key  of\nthe ",(0,i.kt)("inlineCode",{parentName:"p"},"admin")," user. The transaction timeout is set to 20 seconds. This means that the SDK would wait for 20 seconds to receive the\ntransaction's status and receipt synchronously. Once timeout happens, the SDK needs to pool for the transaction status asynchronously."),(0,i.kt)("p",null,"The ",(0,i.kt)("inlineCode",{parentName:"p"},"Session()")," would return the ",(0,i.kt)("inlineCode",{parentName:"p"},"DBSession")," implementation that allows the user to execute various database transactions and queries.\nThe ",(0,i.kt)("inlineCode",{parentName:"p"},"DBSession")," implementation supports the following methods:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"// DBSession captures user's session\ntype DBSession interface {\n    // DBsTx starts a Database Administration Transaction\n    DBsTx() (DBsTxContext, error)\n    // UserTx starts a User Administration Transaction\n    UsersTx() (UsersTxContext, error)\n    // DataTx starts a Data Transaction\n    DataTx(options ...TxContextOption) (DataTxContext, error)\n    // LoadDataTx loads a pre-compileted data transaction\n    LoadDataTx(*types.DataTxEnvelope) (LoadedDataTxContext, error)\n    // ConfigTx starts a Cluster Configuration Transaction\n    ConfigTx() (ConfigTxContext, error)\n    // Provenance returns a provenance querier that supports various provenance queries\n    Provenance() (Provenance, error)\n    // Ledger returns a ledger querier that supports various ledger queries\n    Ledger() (Ledger, error)\n    // JSONQuery returns a JSON querier that supports complex queries on value fields using JSON syntax\n    JSONQuery() (JSONQuery, error)\n}\n")),(0,i.kt)("p",null,"Once the user gets the ",(0,i.kt)("inlineCode",{parentName:"p"},"DBSession"),", any types of transaction can be started"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"    tx, err := session.DBsTx()\n")))}d.isMDXComponent=!0}}]);