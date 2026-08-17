import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "@clerk/react";

export default function ProtectedRoute() {
    const { isLoaded, isSignedIn } = useAuth();

    if (!isLoaded) {
        return <div className="auth-loading">Loading...</div>;
    }

    if (!isSignedIn) {
        return <Navigate to="/auth" replace />;
    }

    return <Outlet />;
}
