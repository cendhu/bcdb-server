"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[3149],{3905:function(e,t,n){n.d(t,{Zo:function(){return p},kt:function(){return f}});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var l=r.createContext({}),u=function(e){var t=r.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},p=function(e){var t=u(e.components);return r.createElement(l.Provider,{value:t},e.children)},s={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},d=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,i=e.originalType,l=e.parentName,p=c(e,["components","mdxType","originalType","parentName"]),d=u(n),f=o,g=d["".concat(l,".").concat(f)]||d[f]||s[f]||i;return n?r.createElement(g,a(a({ref:t},p),{},{components:n})):r.createElement(g,a({ref:t},p))}));function f(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=n.length,a=new Array(i);a[0]=d;var c={};for(var l in t)hasOwnProperty.call(t,l)&&(c[l]=t[l]);c.originalType=e,c.mdxType="string"==typeof e?e:o,a[1]=c;for(var u=2;u<i;u++)a[u]=n[u];return r.createElement.apply(null,a)}return r.createElement.apply(null,n)}d.displayName="MDXCreateElement"},5823:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return c},contentTitle:function(){return l},metadata:function(){return u},toc:function(){return p},default:function(){return d}});var r=n(7462),o=n(3366),i=(n(7294),n(3905)),a=["components"],c={id:"overview",title:"Overview"},l=void 0,u={unversionedId:"getting-started/launching-one-node/overview",id:"getting-started/launching-one-node/overview",isDocsHomePage:!1,title:"Overview",description:"\x3c!--",source:"@site/docs/getting-started/launching-one-node/overview.md",sourceDirName:"getting-started/launching-one-node",slug:"/getting-started/launching-one-node/overview",permalink:"/orion-server/docs/getting-started/launching-one-node/overview",tags:[],version:"current",frontMatter:{id:"overview",title:"Overview"},sidebar:"Documentation",previous:{title:"Using Orion",permalink:"/orion-server/docs/getting-started/guide"},next:{title:"Orion Executable",permalink:"/orion-server/docs/getting-started/launching-one-node/binary"}},p=[],s={toc:p};function d(e){var t=e.components,n=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,r.Z)({},s,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("p",null,"A single Orion node can be started using one of the following two methods:"),(0,i.kt)("ol",null,(0,i.kt)("li",{parentName:"ol"},(0,i.kt)("inlineCode",{parentName:"li"},"orion")," binary executable - ",(0,i.kt)("a",{parentName:"li",href:"binary"},"Click Here")),(0,i.kt)("li",{parentName:"ol"},"docker container - ",(0,i.kt)("a",{parentName:"li",href:"docker"},"Click Here"))),(0,i.kt)("p",null,"For a simple setup, we can use the default configuration and crypto materials supplied in the ",(0,i.kt)("inlineCode",{parentName:"p"},"deployment/")," directory of Hyperledger Orion repository."))}d.isMDXComponent=!0}}]);