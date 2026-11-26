import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { ToastContext, type ToastApi } from "@/context/toast-context";
import DocumentsPage from "@/pages/DocumentsPage";

import type { Document } from "@/types/document";

const documentService = vi.hoisted(() => ({
    getDocuments: vi.fn(),
    searchDocuments: vi.fn(),
    deleteDocument: vi.fn(),
}));

vi.mock("@/services/documentService", () => documentService);

// The upload component is covered by its own concerns; here it is only a
// button that reports a successful upload so the page's reload can be tested.
vi.mock("@/components/documents/DocumentUpload", () => ({
    default: ({ onUploadSuccess }: { onUploadSuccess: () => void }) => (
        <button type="button" onClick={onUploadSuccess}>
            Fake upload done
        </button>
    ),
}));

const toast: ToastApi = {
    success: vi.fn(),
    error: vi.fn(),
    info: vi.fn(),
    warning: vi.fn(),
};

const documents: Document[] = [
    { id: 1, name: "guide.md", type: "md", content: "guide body", created_at: "2026-01-01T00:00:00Z" },
    { id: 2, name: "notes.txt", type: "txt", content: "notes body", created_at: "2026-01-02T00:00:00Z" },
];

function renderPage() {
    return render(
        <ToastContext.Provider value={toast}>
            <DocumentsPage />
        </ToastContext.Provider>,
    );
}

describe("DocumentsPage", () => {

    beforeEach(() => {
        vi.clearAllMocks();
        vi.spyOn(console, "error").mockImplementation(() => {});
    });

    it("shows a loading state and then the documents", async () => {
        documentService.getDocuments.mockResolvedValue(documents);

        renderPage();

        expect(screen.getByText("Loading documents")).toBeInTheDocument();
        expect(await screen.findByText("guide.md")).toBeInTheDocument();
        expect(screen.getByText("notes.txt")).toBeInTheDocument();
        expect(screen.queryByText("Loading documents")).not.toBeInTheDocument();
    });

    it("shows the empty state when there are no documents", async () => {
        documentService.getDocuments.mockResolvedValue([]);

        renderPage();

        expect(await screen.findByText("No documents yet")).toBeInTheDocument();
    });

    it("shows an error and loads again when 'Try again' is clicked", async () => {
        documentService.getDocuments
            .mockRejectedValueOnce(new Error("boom"))
            .mockResolvedValueOnce(documents);

        renderPage();

        expect(await screen.findByText("Failed to load documents.")).toBeInTheDocument();

        await userEvent.click(screen.getByRole("button", { name: "Try again" }));

        expect(await screen.findByText("guide.md")).toBeInTheDocument();
        expect(screen.queryByText("Failed to load documents.")).not.toBeInTheDocument();
        expect(documentService.getDocuments).toHaveBeenCalledTimes(2);
    });

    it("searches and shows only the matching documents", async () => {
        documentService.getDocuments.mockResolvedValue(documents);
        documentService.searchDocuments.mockResolvedValue([documents[0]]);

        renderPage();
        await screen.findByText("guide.md");

        await userEvent.type(screen.getByPlaceholderText("Search documents..."), "g");

        await vi.waitFor(() => {
            expect(screen.queryByText("notes.txt")).not.toBeInTheDocument();
        });
        expect(screen.getByText("guide.md")).toBeInTheDocument();
        expect(documentService.searchDocuments).toHaveBeenCalledWith("g");
    });

    it("shows a 'no documents found' state when the search has no results", async () => {
        documentService.getDocuments.mockResolvedValue(documents);
        documentService.searchDocuments.mockResolvedValue([]);

        renderPage();
        await screen.findByText("guide.md");

        await userEvent.type(screen.getByPlaceholderText("Search documents..."), "zzz");

        expect(await screen.findByText("No documents found")).toBeInTheDocument();
    });

    it("reloads the full list when the search is cleared", async () => {
        documentService.getDocuments.mockResolvedValue(documents);
        documentService.searchDocuments.mockResolvedValue([documents[0]]);

        renderPage();
        await screen.findByText("guide.md");

        const input = screen.getByPlaceholderText("Search documents...");
        await userEvent.type(input, "g");
        await vi.waitFor(() => {
            expect(screen.queryByText("notes.txt")).not.toBeInTheDocument();
        });

        await userEvent.clear(input);

        expect(await screen.findByText("notes.txt")).toBeInTheDocument();
        expect(documentService.getDocuments).toHaveBeenCalledTimes(2);
    });

    it("shows an error toast when the search fails", async () => {
        documentService.getDocuments.mockResolvedValue(documents);
        documentService.searchDocuments.mockRejectedValue(new Error("boom"));

        renderPage();
        await screen.findByText("guide.md");

        await userEvent.type(screen.getByPlaceholderText("Search documents..."), "g");

        await vi.waitFor(() => {
            expect(toast.error).toHaveBeenCalledWith("Failed to search documents.");
        });
    });

    it("asks for confirmation before deleting and removes the document", async () => {
        documentService.getDocuments.mockResolvedValue(documents);
        documentService.deleteDocument.mockResolvedValue(undefined);

        renderPage();
        await screen.findByText("guide.md");

        await userEvent.click(screen.getByRole("button", { name: "Delete guide.md" }));

        expect(screen.getByRole("alertdialog")).toBeInTheDocument();
        expect(documentService.deleteDocument).not.toHaveBeenCalled();

        await userEvent.click(screen.getByRole("button", { name: "Delete" }));

        await vi.waitFor(() => {
            expect(screen.queryByText("guide.md")).not.toBeInTheDocument();
        });
        expect(documentService.deleteDocument).toHaveBeenCalledWith(1);
        expect(toast.success).toHaveBeenCalledWith("Document deleted.");
    });

    it("does not delete when the confirmation is cancelled", async () => {
        documentService.getDocuments.mockResolvedValue(documents);

        renderPage();
        await screen.findByText("guide.md");

        await userEvent.click(screen.getByRole("button", { name: "Delete guide.md" }));
        await userEvent.click(screen.getByRole("button", { name: "Cancel" }));

        expect(documentService.deleteDocument).not.toHaveBeenCalled();
        expect(screen.getByText("guide.md")).toBeInTheDocument();
    });

    it("reloads the list after a successful upload", async () => {
        documentService.getDocuments
            .mockResolvedValueOnce([documents[0]])
            .mockResolvedValueOnce(documents);

        renderPage();
        await screen.findByText("guide.md");
        expect(screen.queryByText("notes.txt")).not.toBeInTheDocument();

        await userEvent.click(screen.getByRole("button", { name: "Fake upload done" }));

        expect(await screen.findByText("notes.txt")).toBeInTheDocument();
        expect(documentService.getDocuments).toHaveBeenCalledTimes(2);
    });
});
