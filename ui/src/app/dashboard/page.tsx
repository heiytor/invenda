"use client";

export default function Dashboard() {
    return (
        <div className="dashboard">
            <aside className="sidebar">
                <nav>
                    <ul>
                        <li><a href="#home">Home</a></li>
                        <li><a href="#about">About</a></li>
                        <li><a href="#services">Services</a></li>
                        <li><a href="#contact">Contact</a></li>
                    </ul>
                </nav>
            </aside>
            <main className="content">
                <h1>Welcome to the Dashboard</h1>
                {/* Add content for your dashboard here */}
            </main>

            <style jsx>{`
                .dashboard {
                    display: flex;
                    height: 100vh;
                }

                .sidebar {
                    width: 250px;
                    background-color: #2c3e50;
                    color: white;
                    padding: 20px;
                }

                .sidebar nav ul {
                    list-style: none;
                    padding: 0;
                }

                .sidebar nav ul li {
                    margin-bottom: 10px;
                }

                .sidebar nav ul li a {
                    color: white;
                    text-decoration: none;
                    display: block;
                    padding: 10px;
                    border-radius: 4px;
                }

                .sidebar nav ul li a:hover {
                    background-color: #34495e;
                }

                .content {
                    flex-grow: 1;
                    padding: 20px;
                }
            `}</style>
        </div>
    );
}

