import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { ToastContext, type ToastApi } from "@/context/toast-context";
import NotesPage from "@/pages/NotesPage";

import type { Note } from "@/types/note";

const noteService = vi.hoisted(() => ({
    getNotes: vi.fn(),
    createNote: vi.fn(),
    updateNote: vi.fn(),
    deleteNote: vi.fn(),
}));

vi.mock("@/services/noteService", () => noteService);

const toast: ToastApi = {
    success: vi.fn(),
    error: vi.fn(),
    info: vi.fn(),
    warning: vi.fn(),
};

const notes: Note[] = [
    { id: 1, title: "Alpha", content: "first body" },
    { id: 2, title: "Beta", content: "second body" },
];

function renderPage() {
    return render(
        <ToastContext.Provider value={toast}>
            <NotesPage />
        </ToastContext.Provider>,
    );
}

describe("NotesPage", () => {

    beforeEach(() => {
        vi.clearAllMocks();
        vi.spyOn(console, "error").mockImplementation(() => {});
    });

    it("shows a loading state and then the notes", async () => {
        noteService.getNotes.mockResolvedValue(notes);

        renderPage();

        expect(screen.getByText("Loading notes")).toBeInTheDocument();
        expect(await screen.findByText("Alpha")).toBeInTheDocument();
        expect(screen.getByText("Beta")).toBeInTheDocument();
        expect(screen.queryByText("Loading notes")).not.toBeInTheDocument();
    });

    it("shows the empty state when there are no notes", async () => {
        noteService.getNotes.mockResolvedValue([]);

        renderPage();

        expect(await screen.findByText("No notes yet")).toBeInTheDocument();
    });

    it("shows an error and loads again when 'Try again' is clicked", async () => {
        noteService.getNotes
            .mockRejectedValueOnce(new Error("boom"))
            .mockResolvedValueOnce(notes);

        renderPage();

        expect(await screen.findByText("Failed to load notes.")).toBeInTheDocument();

        await userEvent.click(screen.getByRole("button", { name: "Try again" }));

        expect(await screen.findByText("Alpha")).toBeInTheDocument();
        expect(screen.queryByText("Failed to load notes.")).not.toBeInTheDocument();
        expect(noteService.getNotes).toHaveBeenCalledTimes(2);
    });

    it("creates a note and shows it in the list", async () => {
        noteService.getNotes.mockResolvedValue(notes);
        noteService.createNote.mockResolvedValue({
            id: 3,
            title: "Gamma",
            content: "new body",
        });

        renderPage();
        await screen.findByText("Alpha");

        await userEvent.type(screen.getByPlaceholderText("Title"), "Gamma");
        await userEvent.type(screen.getByPlaceholderText("Content"), "new body");
        await userEvent.click(screen.getByRole("button", { name: "Save Note" }));

        expect(await screen.findByText("new body")).toBeInTheDocument();
        expect(noteService.createNote).toHaveBeenCalledWith({
            title: "Gamma",
            content: "new body",
        });
        expect(toast.success).toHaveBeenCalledWith("Note created.");
    });

    it("reports an error toast and keeps the form when creating fails", async () => {
        noteService.getNotes.mockResolvedValue(notes);
        noteService.createNote.mockRejectedValue(new Error("boom"));

        renderPage();
        await screen.findByText("Alpha");

        await userEvent.type(screen.getByPlaceholderText("Title"), "Gamma");
        await userEvent.type(screen.getByPlaceholderText("Content"), "x");
        await userEvent.click(screen.getByRole("button", { name: "Save Note" }));

        await vi.waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith("Failed to create note.");
        });
        expect(screen.getByPlaceholderText("Title")).toHaveValue("Gamma");
    });

    it("asks for confirmation before deleting and removes the note", async () => {
        noteService.getNotes.mockResolvedValue(notes);
        noteService.deleteNote.mockResolvedValue(undefined);

        renderPage();
        await screen.findByText("Alpha");

        await userEvent.click(screen.getByRole("button", { name: "Delete Alpha" }));

        expect(screen.getByRole("alertdialog")).toBeInTheDocument();
        expect(noteService.deleteNote).not.toHaveBeenCalled();

        await userEvent.click(screen.getByRole("button", { name: "Delete" }));

        await vi.waitFor(() => {
            expect(screen.queryByText("first body")).not.toBeInTheDocument();
        });
        expect(noteService.deleteNote).toHaveBeenCalledWith(1);
        expect(toast.success).toHaveBeenCalledWith("Note deleted.");
    });

    it("does not delete when the confirmation is cancelled", async () => {
        noteService.getNotes.mockResolvedValue(notes);

        renderPage();
        await screen.findByText("Alpha");

        await userEvent.click(screen.getByRole("button", { name: "Delete Alpha" }));
        await userEvent.click(screen.getByRole("button", { name: "Cancel" }));

        expect(noteService.deleteNote).not.toHaveBeenCalled();
        expect(screen.getByText("first body")).toBeInTheDocument();
    });

    it("loads a note into the form to edit and saves the update", async () => {
        noteService.getNotes.mockResolvedValue(notes);
        noteService.updateNote.mockResolvedValue({
            id: 1,
            title: "Alpha v2",
            content: "first body",
        });

        renderPage();
        await screen.findByText("Alpha");

        await userEvent.click(screen.getByRole("button", { name: "Edit Alpha" }));

        expect(screen.getByPlaceholderText("Title")).toHaveValue("Alpha");

        await userEvent.clear(screen.getByPlaceholderText("Title"));
        await userEvent.type(screen.getByPlaceholderText("Title"), "Alpha v2");
        await userEvent.click(screen.getByRole("button", { name: "Update Note" }));

        expect(await screen.findByRole("heading", { name: "Alpha v2" })).toBeInTheDocument();
        expect(noteService.updateNote).toHaveBeenCalledWith(1, {
            title: "Alpha v2",
            content: "first body",
        });
        expect(screen.getByRole("button", { name: "Save Note" })).toBeInTheDocument();
    });
});
