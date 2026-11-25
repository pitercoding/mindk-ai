import { useCallback, useMemo, useState, type ReactNode } from "react";

import {
    AlertIcon,
    CheckIcon,
    CloseIcon,
    InfoIcon,
} from "@/components/common/icons";

import {
    ToastContext,
    type ToastApi,
    type ToastKind,
} from "@/context/toast-context";

interface Toast {
    id: number;
    kind: ToastKind;
    message: string;
}

const AUTO_DISMISS_MS: Record<ToastKind, number> = {
    success: 3500,
    info: 4000,
    warning: 5000,
    error: 6000,
};

const MAX_VISIBLE = 4;

let nextToastId = 1;

function ToastIcon({ kind }: { kind: ToastKind }) {

    if (kind === "success") {
        return <CheckIcon size={16} />;
    }

    if (kind === "info") {
        return <InfoIcon size={16} />;
    }

    return <AlertIcon size={16} />;
}

export default function ToastProvider({
    children,
}: {
    children: ReactNode;
}) {

    const [toasts, setToasts] = useState<Toast[]>([]);

    const dismiss = useCallback((id: number) => {
        setToasts((current) => current.filter((toast) => toast.id !== id));
    }, []);

    const show = useCallback((kind: ToastKind, message: string) => {

        const id = nextToastId++;

        setToasts((current) =>
            // Ignore an identical toast that is already on screen.
            current.some((toast) => toast.kind === kind && toast.message === message)
                ? current
                : [...current, { id, kind, message }].slice(-MAX_VISIBLE),
        );

        window.setTimeout(() => dismiss(id), AUTO_DISMISS_MS[kind]);

    }, [dismiss]);

    const api = useMemo<ToastApi>(() => ({
        success: (message) => show("success", message),
        error: (message) => show("error", message),
        info: (message) => show("info", message),
        warning: (message) => show("warning", message),
    }), [show]);

    return (
        <ToastContext.Provider value={api}>

            {children}

            <div
                className="toast-viewport"
                role="region"
                aria-label="Notifications"
            >
                {toasts.map((toast) => (
                    <div
                        key={toast.id}
                        className={`toast toast-${toast.kind}`}
                        role={toast.kind === "error" ? "alert" : "status"}
                    >
                        <span className="toast-icon">
                            <ToastIcon kind={toast.kind} />
                        </span>

                        <p className="toast-message">
                            {toast.message}
                        </p>

                        <button
                            type="button"
                            className="btn btn-ghost btn-icon toast-close"
                            aria-label="Dismiss notification"
                            onClick={() => dismiss(toast.id)}
                        >
                            <CloseIcon size={14} />
                        </button>
                    </div>
                ))}
            </div>

        </ToastContext.Provider>
    );
}
