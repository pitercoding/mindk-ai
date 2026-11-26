import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import NoteForm from "@/components/notes/NoteForm";

function setup(
    props: Partial<Parameters<typeof NoteForm>[0]> = {},
) {
    const onCreated = vi.fn().mockResolvedValue(true);
    const onUpdated = vi.fn().mockResolvedValue(true);

    render(
        <NoteForm
            isSaving={false}
            onCreated={onCreated}
            onUpdated={onUpdated}
            {...props}
        />,
    );

    return { onCreated, onUpdated };
}

describe("NoteForm", () => {

    it("creates a note with trimmed values and clears the form", async () => {
        const { onCreated } = setup();

        await userEvent.type(screen.getByPlaceholderText("Title"), "  My title  ");
        await userEvent.type(screen.getByPlaceholderText("Content"), "  Some content ");
        await userEvent.click(screen.getByRole("button", { name: "Save Note" }));

        expect(onCreated).toHaveBeenCalledWith({
            title: "My title",
            content: "Some content",
        });
        expect(screen.getByPlaceholderText("Title")).toHaveValue("");
        expect(screen.getByPlaceholderText("Content")).toHaveValue("");
    });

    it("does not submit when the title or content is blank", async () => {
        const { onCreated } = setup();

        await userEvent.type(screen.getByPlaceholderText("Title"), "Only title");
        await userEvent.click(screen.getByRole("button", { name: "Save Note" }));

        expect(onCreated).not.toHaveBeenCalled();
    });

    it("keeps what was typed when saving fails", async () => {
        const onCreated = vi.fn().mockResolvedValue(false);
        setup({ onCreated });

        await userEvent.type(screen.getByPlaceholderText("Title"), "T");
        await userEvent.type(screen.getByPlaceholderText("Content"), "C");
        await userEvent.click(screen.getByRole("button", { name: "Save Note" }));

        expect(screen.getByPlaceholderText("Title")).toHaveValue("T");
        expect(screen.getByPlaceholderText("Content")).toHaveValue("C");
    });

    it("fills the form and updates the selected note", async () => {
        const { onCreated, onUpdated } = setup({
            selectedNote: { id: 7, title: "Old", content: "Old body" },
        });

        expect(screen.getByPlaceholderText("Title")).toHaveValue("Old");
        expect(screen.getByPlaceholderText("Content")).toHaveValue("Old body");

        await userEvent.clear(screen.getByPlaceholderText("Title"));
        await userEvent.type(screen.getByPlaceholderText("Title"), "New");
        await userEvent.click(screen.getByRole("button", { name: "Update Note" }));

        expect(onUpdated).toHaveBeenCalledWith(7, {
            title: "New",
            content: "Old body",
        });
        expect(onCreated).not.toHaveBeenCalled();
    });

    it("disables the button and shows progress while saving", () => {
        setup({ isSaving: true });

        expect(screen.getByRole("button", { name: "Saving..." })).toBeDisabled();
    });
});
