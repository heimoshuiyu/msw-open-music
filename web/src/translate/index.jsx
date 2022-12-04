import * as React from 'react';
import { createContext } from "react";
import MAP_zh_CN from "./zh_CN";

const LANG_OPTIONS = {
  "en-US": {
    name: "English",
    langMap: {},
    matches: ["en-US", "en"],
  },
  "zh-CN": {
    name: "中文（简体）",
    langMap: MAP_zh_CN,
    matches: ["zh-CN", "zh"],
  },
};

const langCodeContext = createContext("en-US");

function tr(text, langCode) {
  const option = LANG_OPTIONS[langCode];
  if (option === undefined) {
    return text;
  }
  const langMap = LANG_OPTIONS[langCode].langMap;

  const translatedText = langMap[text.toLowerCase()];
  if (translatedText === undefined) {
    return text;
  }

  return translatedText;
}

function Tr(text) {
  return (
    <langCodeContext.Consumer>
      {({ langCode }) => {
        return tr(text, langCode);
      }}
    </langCodeContext.Consumer>
  );
}

export { tr, Tr, LANG_OPTIONS, langCodeContext };
