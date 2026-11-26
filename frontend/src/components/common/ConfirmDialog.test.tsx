import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import ConfirmDialog from "@/components/common/ConfirmDialog";

function setup(props: Partial<Parameters<typeof ConfirmDialog>[0]> = {}) {
    const onConfirm = vi.fn();
    const onCancel = vi.fn();

    render(
        <ConfirmDialog
            title="Delete note?"
            message="It will be gone."
            onConfirm={onConfirm}
            onCancel={onCancel}
            {...props}
        />,
    );

    return { onConfirm, onCancel };
}

describe("ConfirmDialog", () => {

    it("shows the title, message and default button labels", () => {
        setup();

        expect(screen.getByRole("alertdialog")).toBeInTheDocument();
        expect(screen.getByText("Delete note?")).toBeInTheDocument();
        expect(screen.getByText("It will be gone.")).toBeInTheDocument();
        expect(screen.getByRole("button", { name: "Delete" })).toBeInTheDocument();
        expect(screen.getByRole("button", { name: "Cancel" })).toBeInTheDocument();
    });

    it("uses custom button labels", () => {
        setup({ confirmLabel: "Remove", cancelLabel: "Keep" });

        expect(screen.getByRole("button", { name: "Remove" })).toBeInTheDocument();
        expect(screen.getByRole("button", { name: "Keep" })).toBeInTheDocument();
    });

    it("focuses the safe (cancel) action first", () => {
        setup();

        expect(screen.getByRole("button", { name: "Cancel" })).toHaveFocus();
    });

    it("calls onConfirm when the confirm button is clicked", async () => {
        const { onConfirm, onCancel } = setup();

        await userEvent.click(screen.getByRole("button", { name: "Delete" }));

        expect(onConfirm).toHaveBeenCalledTimes(1);
        expect(onCancel).not.toHaveBeenCalled();
    });

    it("calls onCancel when the cancel button is clicked", async () => {
        const { onConfirm, onCancel } = setup();

        await userEvent.click(screen.getByRole("button", { name: "Cancel" }));

        expect(onCancel).toHaveBeenCalledTimes(1);
        expect(onConfirm).not.toHaveBeenCalled();
    });

    it("calls onCancel when Escape is pressed", async () => {
        const { onCancel } = setup();

        await userEvent.keyboard("{Escape}");

        expect(onCancel).toHaveBeenCalledTimes(1);
    });
});
