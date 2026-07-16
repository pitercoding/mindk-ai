import { Outlet } from "react-router-dom";

import Sidebar from "@/components/layout/Sidebar";
import Header from "@/components/layout/Header";

export default function DashboardLayout() {
    return (
        <div className="dashboard-layout">

            <Sidebar />

            <div className="dashboard-content">

                <Header />

                <main>
                    <Outlet />
                </main>

            </div>

        </div>
    );
}
