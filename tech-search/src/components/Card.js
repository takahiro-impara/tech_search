import { Card } from 'semantic-ui-react';

const CardProps = ({Title, Date, Url, Company}) => {
  const company_bucket = process.env.REACT_APP_BACKETENDPOINT + Company + process.env.REACT_APP_IMEXTENSION
  return (
    <Card
      href={Url}
      image={company_bucket}
      header={Company}
      meta={Date}
      description={Title}
    />
  )
}

export default CardProps