/**
 * @author n1474335 [n1474335@gmail.com]
 * @copyright Crown Copyright 2016
 * @license Apache-2.0
 */

import Operation from "../Operation.mjs";
import cptable from "codepage";
import { CHR_ENC_CODE_PAGES } from "../lib/ChrEnc.mjs";
import iconv from 'iconv-lite'

/**
 * Decode text operation
 */
class DecodeText extends Operation {

  /**
   * DecodeText constructor
   */
  constructor() {
    super();

    this.name = "Decode text";
    this.module = "Encodings";
    this.description = [
      "Decodes text from the chosen character encoding.",
      "<br><br>",
      "Supported charsets are:",
      "<ul>",
      Object.keys(CHR_ENC_CODE_PAGES).map(e => `<li>${e}</li>`).join("\n"),
      "</ul>",
    ].join("\n");
    this.infoURL = "https://wikipedia.org/wiki/Character_encoding";
    this.inputType = "ArrayBuffer";
    this.outputType = "string";
    this.args = [
      {
        "name": "Encoding",
        "type": "option",
        "value": Object.keys(CHR_ENC_CODE_PAGES)
      }
    ];
  }

  /**
   * @param {ArrayBuffer} input
   * @param {Object[]} args
   * @returns {string}
   */
  run(input, args) {
    if (typeof input == 'string') {
      const encoder = new TextEncoder()
      const encoded = encoder.encode(input)
      const format = CHR_ENC_CODE_PAGES[args[0]];
      let str = ''
      try {
        str = cptable.utils.decode(format, encoded)
      }
      catch {
        str = input
      }
      return str;
    } else {
      const format = CHR_ENC_CODE_PAGES[args[0]];
      return cptable.utils.decode(format, new Uint8Array(input));
    }
  }

}

export default DecodeText;
