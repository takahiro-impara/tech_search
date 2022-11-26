const Card = ({Title, Date, Url, Company}) => {
  const company_bucket = "http://tech-search.s3.ap-northeast-1.amazonaws.com/" + Company + ".png"
  return (
    <div className="card-item">
      <div className="card-img">
        <img alt={Company} src={company_bucket} />
      </div>
      <div className="card-title">
        <a href={Url}><h3>{Title}</h3></a>
      </div>
      <div className="card-subInfo">
        <div className="company">{Company}</div>
        <div className="date">{Date}</div>
      </div>
    </div>
  )
}

export default Card