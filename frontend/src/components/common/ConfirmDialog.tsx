import { useEffect, useRef } from "react";
import { createPortal } from "react-dom";

interface ConfirmDialogProps {
    title: string;
    message: string;
    confirmLabel?: string;
    cancelLabel?: string;
    onConfirm: () => void;
    onCancel: () => void;
}

export default function ConfirmDialog({
    title,
    message,
    confirmLabel = "Delete",
    cancelLabel = "Cancel",
    onConfirm,
    onCancel,
}: ConfirmDialogProps) {

    const cancelRef = useRef<HTMLButtonElement | null>(null);
    const onCancelRef = useRef(onCancel);

    useEffect(() => {
        onCancelRef.current = onCancel;
    });

    // Focus the safe action first and let Escape dismiss the dialog.
    useEffect(() => {

        cancelRef.current?.focus();

        function handleKeyDown(event: KeyboardEvent) {
            if (event.key === "Escape") {
                onCancelRef.current();
            }
        }

        window.addEventListener("keydown", handleKeyDown);

        return () => {
            window.removeEventListener("keydown", handleKeyDown);
        };

    }, []);

    return createPortal(

        <div
            className="confirm-backdrop"
            onMouseDown={onCancel}
        >

            <div
                className="confirm-dialog"
                role="alertdialog"
                aria-modal="true"
                aria-labelledby="confirm-dialog-title"
                aria-describedby="confirm-dialog-message"
                onMouseDown={(event) => event.stopPropagation()}
            >

                <h2 id="confirm-dialog-title">
                    {title}
                </h2>

                <p id="confirm-dialog-message">
                    {message}
                </p>

                <div className="confirm-actions">

                    <button
                        ref={cancelRef}
                        type="button"
                        className="confirm-button"
                        onClick={onCancel}
                    >
                        {cancelLabel}
                    </button>

                    <button
                        type="button"
                        className="confirm-button confirm-button-danger"
                        onClick={onConfirm}
                    >
                        {confirmLabel}
                    </button>

                </div>

            </div>

        </div>,

        document.body,
    );
}
