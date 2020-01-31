import React from 'react'

const Customer = ({ clients }) => {
  return (
    <div>
      <center><h1>Lista de clientes</h1></center>
      {clients.map((cli) => (
        <div className="card">
          <div className="card-body">
            <h5 className="card-title">{cli.FirstName}</h5>
            <h6 className="card-subtitle mb-2 text-muted">{cli.LastName}</h6>
            <p className="card-text">{cli.Score.toString()}</p>
          </div>
        </div>
        
      ))}
    </div>
  )
};

export default Customer