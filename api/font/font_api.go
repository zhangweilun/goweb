package font

import (

	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"

)

/**
*
* @author willian
* @created 2017-01-28 13:02
* @email 18702515157@163.com
**/

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}


