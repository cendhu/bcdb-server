"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[8214],{3905:function(e,t,n){n.d(t,{Zo:function(){return u},kt:function(){return h}});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var c=r.createContext({}),l=function(e){var t=r.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},u=function(e){var t=l(e.components);return r.createElement(c.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},p=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,c=e.parentName,u=s(e,["components","mdxType","originalType","parentName"]),p=l(n),h=a,f=p["".concat(c,".").concat(h)]||p[h]||d[h]||o;return n?r.createElement(f,i(i({ref:t},u),{},{components:n})):r.createElement(f,i({ref:t},u))}));function h(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,i=new Array(o);i[0]=p;var s={};for(var c in t)hasOwnProperty.call(t,c)&&(s[c]=t[c]);s.originalType=e,s.mdxType="string"==typeof e?e:a,i[1]=s;for(var l=2;l<o;l++)i[l]=n[l];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}p.displayName="MDXCreateElement"},3527:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return s},contentTitle:function(){return c},metadata:function(){return l},toc:function(){return u},default:function(){return p}});var r=n(7462),a=n(3366),o=(n(7294),n(3905)),i=["components"],s={id:"introduction",title:"Hyperledger Orion"},c=void 0,l={unversionedId:"external/introduction",id:"external/introduction",isDocsHomePage:!1,title:"Hyperledger Orion",description:"\x3c!--",source:"@site/docs/external/introduction.md",sourceDirName:"external",slug:"/external/introduction",permalink:"/orion-server/docs/external/introduction",tags:[],version:"current",frontMatter:{id:"introduction",title:"Hyperledger Orion"},sidebar:"Documentation",next:{title:"Using Orion",permalink:"/orion-server/docs/external/getting-started/guide"}},u=[{value:"High Level Architecture",id:"high-level-architecture",children:[],level:2}],d={toc:u};function p(e){var t=e.components,n=(0,a.Z)(e,i);return(0,o.kt)("wrapper",(0,r.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"Orion is a ",(0,o.kt)("strong",{parentName:"p"},"key-value/document database")," with certain blockchain properties such as"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Tamper Evident"),": Data cannot be tampered with, without it going unnoticed. At any point in time, a user can request the database to provide proof for an existance of a transaction or data, and verify the same to ensure data integrity."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Non-Repudiation"),": A user who submitted a transaction to make changes to data cannot deny submitting the transaction later."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Crypto-based Authentication"),": A user that submitted a query or transaction is always authenticated using digital signature."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Confidentiality and Access Control"),": Each data item can have an access control list (ACL) to dictate which users can read from it and which users can write to it. Each user needs to authenticate themselves by providing their digital signature to read or write to data. Depending on the access rule defined for data, sometimes more than one users need to authenticate themselves together to read or write to data."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Serialization Isolation Level"),": It ensures a safe and consistent transaction execution."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Provenance Queries"),": All historical changes to the data are maintained separately in a persisted graph data structure so that a user can execute query on those historical changes to understand the lineage of each data item.")),(0,o.kt)("p",null,"OrionB",(0,o.kt)("strong",{parentName:"p"},"DOES NOT")," have the following two blockchain properties:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Smart-Contracts"),": A set of functions that manage data on the blockchain ledger. Transactions are invocations of one or more smart contract's functions."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Decentralization of Trust"),": A permissioned setup of known but untrusted organizations each operating their own independent database nodes but connected together to form a blockchain network. As one node cannot trust the execution results of another node, ordering transaction must be done with a BFT protocol\nand all transactions need to be independently executed on all nodes.")),(0,o.kt)("h2",{id:"high-level-architecture"},"High Level Architecture"),(0,o.kt)("p",null,"Figure 1 presents the high level architecture of Orion."),(0,o.kt)("img",{src:"./docs/figures/high-level-architecture.png",alt:"drawing",width:"800"}),(0,o.kt)("p",null,"OrionBstores and manages the following five data elements:"),(0,o.kt)("ol",null,(0,o.kt)("li",{parentName:"ol"},(0,o.kt)("strong",{parentName:"li"},"Users"),": Storage of users' credentials such as digital certificate and their privileges. Only these users can access the Orion node."),(0,o.kt)("li",{parentName:"ol"},(0,o.kt)("strong",{parentName:"li"},"Key-Value Pairs"),": Storage of all current/active key-value pairs committed by users of the Orion node."),(0,o.kt)("li",{parentName:"ol"},(0,o.kt)("strong",{parentName:"li"},"Historical Key-Value Pairs"),": Storage of all past/inactive key-value pairs using a graph data structure with additional metadata\nsuch as the user who modified the key-value pair, all previous and next values of the key, transactions which have read or written to\nthe key-value pair, etc... It helps to provide a complete data lineage."),(0,o.kt)("li",{parentName:"ol"},(0,o.kt)("strong",{parentName:"li"},"Authenticated Data Structure"),": Storage of Merkle Patricia Tree where leaf node is nothing but a key-value pair. It helps in\ncreating proofs for the existance or non-existance of a key-value pair, and transaction."),(0,o.kt)("li",{parentName:"ol"},(0,o.kt)("strong",{parentName:"li"},"Hash chain of blocks"),": Storage of cryptographically linked blocks, where each block holds a set of transactions submitted\nby the user along with its commit status, summary of state changes in the form of Merkle Patricia's Root hash, etc... It helps in\ncreating a proof for the existance of a block or a transaction.")),(0,o.kt)("p",null,"The users of the Orion can query these five data elements provided that they have the required privileges and\nalso can perform transactions to modify active key-value pairs. When a user submits a transaction, user receives a transaction receipt\nfrom the Orion node after the commit of a block that includes the transaction. The user can then store the receipt locally for performing\nclient side verification of proof of existance of a key-value pair or a transaction or a block."))}p.isMDXComponent=!0}}]);