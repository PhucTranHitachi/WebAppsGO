import { useState } from "react";
import Navbar from "../components/Navbar";
import {useNavigate} from 'react-router-dom'

export default function AddBook() {
    const navigate = useNavigate()
    const [data, setData] = useState({
        name: "",
        description: "",
        status: ""
    })

    const handleInput = (e) => {
        const name = e.target.name
        const value = e.target.value
        setData({...data, [name]: value})
    }

    const handleSubmit = async(e) => {
        e.preventDefault();
        const { name, description, status } = data;
        try {
            const res = await fetch("http://localhost:8000/v1/items", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ name, description, status })
            });
    
            if (!res.ok) {
                throw new Error("Failed to submit");
            }
    
            const body = await res.json();
            console.log(body);
            navigate("/");
        } catch (error) {
            console.error("Error submitting book:", error);
        }
    };

    return (
        <div>
            <Navbar />
            <div className="container">
                <h1 className="text-center">Add Book</h1>
                <form onSubmit={handleSubmit}>
                    <div className="mb-3">
                        <label htmlFor="name" className="form-label">Name</label>
                        <input type="text" className="form-control" name="name" id="name" onChange={handleInput}/>
                    </div>
                    <div className="mb-3">
                        <label htmlFor="description" className="form-label">Description</label>
                        <input type="text" className="form-control" name="description" id="description" onChange={handleInput}/>
                    </div>
                    <div className="mb-3">
                        <label htmlFor="status" className="form-label">Status</label>
                        <input type="text" className="form-control" name="status" id="status" onChange={handleInput}/>
                    </div>
                    <button type="submit" className="btn btn-primary w-100">Submit</button>
                </form>
            </div>
        </div>
    )
}