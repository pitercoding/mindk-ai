import { describe, expect, it } from "vitest";

import {
    countWords,
    formatDate,
    stripTitleHeading,
    toPlainText,
} from "@/utils/text";

describe("toPlainText", () => {

    it("removes headings, emphasis and inline code", () => {
        expect(toPlainText("# Title\n**bold** and *italic* and `code`"))
            .toBe("Title bold and italic and code");
    });

    it("keeps only the visible text of links and images", () => {
        expect(toPlainText("see [docs](https://x.dev) ![logo](a.png)"))
            .toBe("see docs logo");
    });

    it("drops fenced code blocks", () => {
        expect(toPlainText("before\n```ts\nconst a = 1;\n```\nafter"))
            .toBe("before after");
    });

    it("removes list markers, task boxes and blockquotes", () => {
        expect(toPlainText("- [x] done\n- item\n1. first\n> quote"))
            .toBe("done item first quote");
    });

    it("removes table pipes and separator rows", () => {
        expect(toPlainText("| a | b |\n|---|---|\n| 1 | 2 |"))
            .toBe("a b 1 2");
    });

    it("does not mangle snake_case names", () => {
        expect(toPlainText("use my_var_name here")).toBe("use my_var_name here");
    });

    it("returns an empty string for blank input", () => {
        expect(toPlainText("   \n  ")).toBe("");
    });
});

describe("countWords", () => {

    it("counts words separated by any whitespace", () => {
        expect(countWords("one  two\nthree\tfour")).toBe(4);
    });

    it("returns 0 for empty or blank text", () => {
        expect(countWords("")).toBe(0);
        expect(countWords("   ")).toBe(0);
    });
});

describe("formatDate", () => {

    it("formats a valid date with the fixed en-US locale", () => {
        expect(formatDate("2026-03-05T12:00:00Z")).toBe("Mar 05, 2026");
    });

    it("returns an empty string for an invalid date", () => {
        expect(formatDate("not-a-date")).toBe("");
    });
});

describe("stripTitleHeading", () => {

    it("removes a leading heading that repeats the title", () => {
        expect(stripTitleHeading("# My Note\n\nBody", "My Note")).toBe("Body");
    });

    it("keeps content when the heading is different", () => {
        expect(stripTitleHeading("# Other\nBody", "My Note")).toBe("# Other\nBody");
    });

    it("treats regex characters in the title literally", () => {
        expect(stripTitleHeading("# a.b (c)\nBody", "a.b (c)")).toBe("Body");
    });
});
