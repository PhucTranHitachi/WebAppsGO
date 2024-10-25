import Navbar from "../components/Navbar"
import { useEffect, useState } from "react"

export default function Home() {
    const [data, setData] = useState([])
    const callAPI = async () => {
        const res = await fetch('http://localhost:8000/v1/items')
        const body = await res.json()
        console.log(body)
        setData(body)
    }

    const handleDelete = (id) => {
        return async () => {
            const res = await fetch(`http://localhost:8000/v1/items/${id}`, {
                method: 'DELETE'
            })
            const body = await res.json()
            console.log(body)
            callAPI();
        }
    }

    useEffect(() => {
        callAPI();
    }, [])

    return (
        <div>
            <Navbar />
            <div className="container">
                <table className="table">
                    <thead>
                        <tr>
                            <th scope="col">ID</th>
                            <th scope="col">Name</th>
                            <th scope="col">Description</th>
                            <th scope="col">Status</th>
                            <th scope="col">Created At</th>
                            <th scope="col">Updated At</th>
                        </tr>
                    </thead>
                    <tbody>
                        {data.data && data.data.length > 0 ? (
                            data.data.map((item, index) => (
                                <tr key={index}>
                                    <th scope="row">{item.id}</th>
                                    <td>{item.name}</td>
                                    <td>{item.description}</td>
                                    <td>{item.status}</td>
                                    <td>{item.created_at}</td>
                                    <td>{item.updated_at}</td>
                                    <td>
                                        <a href={`/update/${item.id}`} className="btn btn-warning me-3">Update</a>
                                        <button type="button" className="btn btn-danger" onClick={handleDelete(item.id)}>Delete</button>
                                    </td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan="7" className="text-center">No data available</td>
                            </tr>
                        )}
                    </tbody>


                </table>
            </div>
        </div>
    )
}