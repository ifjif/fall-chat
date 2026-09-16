# 数据库的查询准则

在业务可以进行处理的情况下，查询只需要查出数据就可以
不需要关系页面要显示什么,如：

      select id, username,
        case
          when status = 0 then 'pending'
          when status = 1 then 'approved'
        else 'unknown' as status_text

像上面这种，再查询时还关系了页面的显示逻辑
问题就是：如果后续不想用pending这个，那就需要改sql,这是比较危险的操作

这种可以直接在他自己model中进行处理，已go为例：

    go
      type UserModel struct {
        id int64
        username string
        status int8
      }
      var statuMaps = map[int8]string{
        0:"pending",
        1:"approved"
      }
      func (u *UserModel) GetStatusText() string {
        if text, ok := statuMaps[u.status]; ok {
          return text
        }
        return "unknown"
      }

这样，页面的显示逻辑完全由业务管理，而数据库只负责查询。
除非是那种直接走数据库的，不经过业务处理，可以直接在查询时进行转换。

# 数据库的动态唯一键

有时，我们需要某个数据在某个状态下只有一条，而它不在这个转态时，可以有多条。
如，在im系统中一个 好友请求申请:
在pending和approved时，它只能有一条，不能有多条
而在rejected和expired时，它可以有多条。

如表：{
id int64
from_user_id int64
to_user_id int64
status int8
}

这时，我们可以加一个字段 uniqueKey, 它在pending和approved时，就是 'active',其它就是当前时间戳，来释放唯一

如表：{
id int64
from_user_id int64
to_user_id int64
status int8
uniqueKey string
}
这时 可以给 (from_user_id,to_user_id,uniqueKey)设置唯一索引

**但是，这又出现了一个问题**

from_user_id -> to_user_id 这能保证是唯一的，而
to_user_id -> from_user_id 可以再次插入。
如果要保持绝对唯一，这个就不可以了

这时可以再加一个字段为 `pairKey` 它是 (from*user_id,to_user_id)的结合
min(from_user_id,to_user_id)\_max(from_user_id, to_user_id)
然后 (pairKey,uniqueKey)成为唯一索引
生成它的值，可以在业务逻辑进行处理，也可以在数据库层面进行处理（支持生成列的）
如mysql中:
`pairKey varchar(100) generated always as (concat(least(from_user_id,to_user_id),'*',greatest(from_user_id,to_user_id))) stored`
stored表示会存储到磁盘中，默认是virtual。
这个列的值无法主动设置,不能修改，表达式的值必须是确定的。
