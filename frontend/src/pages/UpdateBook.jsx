import { useEffect, useState } from "react";
import Navbar from "../components/Navbar";
import { useParams, useNavigate } from "react-router-dom";

export default function AddBook() {
    const navigate = useNavigate()
    const {id} = useParams()
    const [data, setData] = useState({
        name: "",
        description: "",
        status: ""
    })

    const callAPI = async () => {
        const res = await fetch(`http://localhost:8000/v1/items/${id}`);
        const body = await res.json()
        setData(body)
        console.log(body)
    }

    useEffect(() => {
        callAPI()
    },[])

    const handleInput = (e) => {
        const name = e.target.name
        const value = e.target.value
        setData({...data, [name]: value})
    }

    const handleSubmit = async (e) => {
        e.preventDefault();
        const { name, description, status } = data;
    
        try {
            const res = await fetch(`http://localhost:8000/v1/items/${id}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ name, description, status })
            });
    
            if (!res.ok) {
                throw new Error('Network response was not ok');
            }
    
            const responseBody = await res.json();
            console.log(responseBody);
            navigate("/");
        } catch (error) {
            console.error('Error:', error);
        }
    };
    

    return (
        <div>
            <Navbar />
            <div className="container">
                <h1 className="text-center">Update Book</h1>
                <form onSubmit={handleSubmit}>
                    <div className="mb-3">
                        <label htmlFor="name" className="form-label">Name</label>
                        <input type="text" className="form-control" name="name" id="name" onChange={handleInput} />
                    </div>
                    <div className="mb-3">
                        <label htmlFor="description" className="form-label">Description</label>
                        <input type="text" className="form-control" name="description" id="description" onChange={handleInput} />
                    </div>
                    <div className="mb-3">
                        <label htmlFor="status" className="form-label">Status</label>
                        <input type="text" className="form-control" name="status" id="status" onChange={handleInput} />
                    </div>
                    <button type="submit" className="btn btn-warning w-100">Submit</button>
                </form>
            </div>
        </div>
    )
}