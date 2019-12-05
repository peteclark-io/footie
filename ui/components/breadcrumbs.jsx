import Link from 'next/link'
import {Typography} from '@material-ui/core'
import {ChevronRight} from '@material-ui/icons'

const Breadcrumbs = (props) => {
  const {crumbs} = props;

  if(crumbs.length === 0){
    return (
      <React.Fragment></React.Fragment>
    )
  }

  crumbs[crumbs.length-1].noChevron = true;

  return (
    <section className="breadcrumbs">
      {crumbs.map(crumb => {
        const useLink = !!crumb.href;
        return (
          <span key={crumb.title}>
            <Typography variant="h6">{useLink ? <Link href={crumb.href}>{crumb.title}</Link> : crumb.title} {crumb.noChevron ? undefined : <ChevronRight className='icon'></ChevronRight>}</Typography>
          </span>
        )
      })}

    </section>
  )
};

export default Breadcrumbs
