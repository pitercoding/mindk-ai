import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import NoteList from "@/components/notes/NoteList";

import type { Note } from "@/types/note";

const notes: Note[] = [
    { id: 1, title: "First", content: "# First\n\nhello **world** again" },
    { id: 2, title: "Second", content: "" },
];

function setup(
    props: Partial<Parameters<typeof NoteList>[0]> = {},
) {
    const handlers = {
        onDelete: vi.fn(),
        onEdit: vi.fn(),
        onCreate: vi.fn(),
    };

    render(
        <NoteList
            notes={notes}
            deletingId={null}
            {...handlers}
            {...props}
        />,
    );

    return handlers;
}

describe("NoteList", () => {

    it("shows an empty state with a create action when there are no notes", async () => {
        const { onCreate } = setup({ notes: [] });

        expect(screen.getByText("No notes yet")).toBeInTheDocument();

        await userEvent.click(
            screen.getByRole("button", { name: /write your first note/i }),
        );

        expect(onCreate).toHaveBeenCalledTimes(1);
    });

    it("renders a card per note with a plain-text preview and word count", () => {
        setup();

        expect(screen.getByText("First")).toBeInTheDocument();
        expect(screen.getByText("hello world again")).toBeInTheDocument();
        expect(screen.getByText("3 words")).toBeInTheDocument();
        expect(screen.getByText("No content yet")).toBeInTheDocument();
        expect(screen.getByText("0 words")).toBeInTheDocument();
    });

    it("calls onEdit and onDelete with the right note", async () => {
        const { onEdit, onDelete } = setup();

        await userEvent.click(screen.getByRole("button", { name: "Edit First" }));
        await userEvent.click(screen.getByRole("button", { name: "Delete Second" }));

        expect(onEdit).toHaveBeenCalledWith(notes[0]);
        expect(onDelete).toHaveBeenCalledWith(2);
    });

    it("disables the buttons of the note being deleted", () => {
        setup({ deletingId: 1 });

        expect(screen.getByRole("button", { name: "Edit First" })).toBeDisabled();
        expect(screen.getByRole("button", { name: "Delete First" })).toBeDisabled();
        expect(screen.getByText("Deleting...")).toBeInTheDocument();
        expect(screen.getByRole("button", { name: "Edit Second" })).toBeEnabled();
    });
});
