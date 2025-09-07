/* eslint camelcase: 0 */


import NodeDish from "./NodeDish.mjs";
import { _wrap, help, bake, _explainExcludedFunction } from "./api.mjs";
import { OperationError, DishError, ExcludedOperationError } from "../core/errors/index.mjs";
import {
  UnescapeUnicodeCharacters as core_UnescapeUnicodeCharacters,
  FromHTMLEntity as core_FromHTMLEntity,
  URLDecode as core_URLDecode,
  // URLEncode as core_URLEncode,
  // DecodeText as core_DecodeText,
  // EncodeText as core_EncodeText,
  // ToBase64 as core_ToBase64,
  // FromBase64 as core_FromBase64,
  // ToHex as core_ToHex,
  // FromHex as core_FromHex,
} from "../core/operations/index.mjs";
import Utils from "../core/Utils.mjs";
/**
 * generateChef
 *
 * Creates decapitalised, wrapped ops in chef object for default export.
 */
function generateChef() {
  return {
    "unescapeUnicodeCharacters": _wrap(core_UnescapeUnicodeCharacters),
    "fromHTMLEntity": _wrap(core_FromHTMLEntity),
    "URLDecode": _wrap(core_URLDecode),
    // "DecodeText": _wrap(core_DecodeText),
    // "EncodeText": _wrap(core_EncodeText),
    // "URLEncode": _wrap(core_URLEncode),
    // "FromBase64": _wrap(core_FromBase64),
    // "ToBase64": _wrap(core_ToBase64),
    // "FromHex": _wrap(core_FromHex),
    // "ToHex": _wrap(core_ToHex),
  };
}

const chef = generateChef();
// Add some additional features to chef object.
chef.help = help;
chef.Dish = NodeDish;

const unescapeUnicodeCharacters = chef.unescapeUnicodeCharacters;
const fromHTMLEntity = chef.fromHTMLEntity;
const URLDecode = chef.URLDecode;
// const DecodeText = chef.DecodeText;
// const EncodeText = chef.EncodeText;
// const URLEncode = chef.URLEncode;
// const FromBase64 = chef.FromBase64;
// const ToBase64 = chef.ToBase64;
// const FromHex = chef.FromHex;
// const ToHex = chef.ToHex;

const operations = [
  fromHTMLEntity,
  unescapeUnicodeCharacters,
  URLDecode,
  // URLEncode,
  // DecodeText,
  // EncodeText,
  // FromBase64,
  // ToBase64,
  // FromHex,
  // ToHex,
];

chef.bake = bake;
export default chef;

// Operations as top level exports.
export {
  operations,
  unescapeUnicodeCharacters,
  fromHTMLEntity,
  URLDecode,
  // URLEncode,
  // DecodeText,
  // EncodeText,
  // FromBase64,
  // ToBase64,
  // FromHex,
  // ToHex,
  NodeDish as Dish,
  bake,
  help,
  OperationError,
  ExcludedOperationError,
  DishError,
  Utils
};
