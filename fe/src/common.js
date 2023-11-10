const common = {}

common.SetFormData = function (e, setFormData) {
  setFormData(prev => {
    return {
      ...prev,
      [e.target.name]: e.target.value
    }
  })
}

export default function Common() {
  return common
}